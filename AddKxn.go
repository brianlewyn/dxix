package dxix

import (
	"strconv"
	"strings"
)

// Todo: [ dy: AddKxn ] ============================================
// * My GO Module: 02

// Sum all similar Kx^n, with the same x^n.
func AddKxn(w, g *string) error {

	// Are there empty input fields?
	err := existsXY(*w, *g)
	if err != nil {
		return err
	}

	// Was the math func written correctly?
	err = correctlyXY(*w, g)
	if err != nil {
		return err
	}

	// Split on the expression Kx^n.
	sKxn := splitY(*g)

	// Finish rebuilding Kxn in its most basic expression.
	for i, Kxn := range sKxn {
		if strings.Contains(Kxn, *w) {

			// Add the + sign when it is only x
			moreLess := strings.HasPrefix(Kxn, "+") || strings.HasPrefix(Kxn, "-")
			if !moreLess {
				sKxn[i] = "+" + Kxn
			}

			// All x or kx => kx^1
			finishBuildingKxn(*w, &sKxn[i])

		} else {
			// All k => kx^0
			sKxn[i] = sKxn[i] + *w + "^0"
		}
	}

	// ? Store k and n in different sets. ==============================
	sN, ssFlt := []float64{}, make([][]float64, len(sKxn))
	for i, kxn := range sKxn {

		// Split kxn into k and n.
		skn := strings.Split(kxn, *w+"^")

		// Convert k to float.
		kFlt, _ := strconv.ParseFloat(skn[0], 64)

		// Convert n to float.
		nFlt, _ := strconv.ParseFloat(skn[1], 64)

		// Save the new k and n.
		sN = append(sN, nFlt)
		ssFlt[i] = []float64{kFlt, nFlt}
	}

	// ? Remove all duplicate n in g. ==============================
	var s []float64
	dicc := make(map[float64]bool)
	for _, flt := range sN {
		if _, key := dicc[flt]; !key {
			dicc[flt] = true
			s = append(s, flt)
		}
	}

	// ? Numbers from highest to lowest, in this case for n. ==============================
	var temp float64
	for x := range s {
		for y := range s {
			if s[x] > s[y] {
				temp = s[x]
				s[x] = s[y]
				s[y] = temp
			}
		}
	}

	sN = s

	// ? Add all the k's with the same n. ==============================
	t1, sskn := 0, make([][]float64, len(sN))
	for _, n := range sN {

		// Store all [k, n], according to its same n.
		t2, sn := 0, make([][]float64, len(sKxn))
		for i := range ssFlt {
			if n == ssFlt[i][1] {
				sn[t2] = append(sn[t2], ssFlt[i][0], ssFlt[i][1])
				t2++
			}
		}

		// Sum all elements k.
		var k float64
		for i := range sn[:t2] {
			k += sn[i][0]
		}

		// Save all the [k, n] new.
		sskn[t1] = append(sskn[t1], k, n)
		t1++
	}

	// ? Rebuild all kxn's and then add. ==============================
	var sknStr string
	for i := range sskn {

		// Convert float to string.
		kStr := strconv.FormatFloat(sskn[i][0], DFrmt, NPrec, 64)
		nStr := strconv.FormatFloat(sskn[i][1], DFrmt, NPrec, 64)

		// Rebuild all the Kxn func.
		kStr = strings.Replace(kStr, Zero, "", 1)
		nStr = strings.Replace(nStr, Zero, "", 1)
		kxn := kStr + *w + "^" + nStr

		// Remove all 1, 0, ^1, ^0
		if kStr == "0" {
			kxn = "0"
		} else if kStr == "1" && nStr == "0" {
			kxn = kStr
		} else if kStr == "1" && nStr == "1" {
			kxn = *w
		} else if kStr != "1" && nStr == "1" {
			kxn = kStr + *w
		} else if kStr != "1" && nStr == "0" {
			kxn = kStr
		} else if kStr == "1" && nStr != "0" && nStr != "1" {
			kxn = *w + "^" + nStr
		}

		// Add space between the sign.
		if i == 0 {
			sknStr += kxn
		} else {
			sknStr += redirectASign(kxn)
		}
	}
	*g = sknStr

	return nil
}
