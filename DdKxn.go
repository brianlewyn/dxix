package dxix

import (
	"strconv"
	"strings"
)

// ! For floating numbers, precision, decimal format and remove .00
var NPrec int = 2
var DFrmt byte = 'f'
var Zero string = ".00"

// Todo: [ -- : Knxm ] ==============================================
// * My custom kxn rule.

// My Struct dkxn.
type dKxn struct {
	k, x, n string
}

// Define the structure of the kxn rule.
func (g *dKxn) knxm() string {

	// [g.k] [g.x] [g.n]
	var kStr, nStr, dg string
	kDot := strings.Contains(g.k, ".")
	nDot := strings.Contains(g.n, ".")

	// Transform the values to the correct data type and then do the calculations
	if !kDot {
		kInt, _ := strconv.Atoi(g.k)
		if !nDot {
			nInt, _ := strconv.Atoi(g.n)
			kStr = strconv.Itoa(kInt * nInt)
			nStr = strconv.Itoa(nInt - 1)
		} else {
			nFlt, _ := strconv.ParseFloat(g.n, 64)
			kStr = strconv.FormatFloat(float64(kInt)*nFlt, DFrmt, NPrec, 64)
			nStr = strconv.FormatFloat(nFlt-1, DFrmt, NPrec, 64)
		}
	} else {
		kFlt, _ := strconv.ParseFloat(g.k, 64)
		if !nDot {
			nInt, _ := strconv.Atoi(g.n)
			kStr = strconv.FormatFloat(kFlt*float64(nInt), DFrmt, NPrec, 64)
			nStr = strconv.Itoa(nInt - 1)
		} else {
			nFlt, _ := strconv.ParseFloat(g.n, 64)
			kStr = strconv.FormatFloat(kFlt*nFlt, DFrmt, NPrec, 64)
			nStr = strconv.FormatFloat(nFlt-1, DFrmt, NPrec, 64)
		}
	}

	// Remove all kxn with k = "0"
	if kStr == "0" {
		dg = "0"
	} else {
		dg = kStr + g.x + "^" + nStr
	}

	return dg
}

// Todo: [ dy: DyKxn ] ============================================
// * My GO Module: 01

// Check if there is any empty input field.
func existsXY(w, g string) error {
	if w == "" && g == "" {
		return myErrGO("Fill in both input fields.")
	} else if w == "" && g != "" {
		return myErrGO("Fill in the field f(a literal).")
	} else if w != "" && g == "" {
		return myErrGO("Fill in the fild = (a math func).")
	} else {
		return nil
	}
}

// Verify if the math func was written correctly.
func correctlyXY(w string, g *string) error {

	// Convert all literals to the same.
	upper, lower := strings.ToUpper(w), strings.ToLower(w)
	if w == upper {
		if !strings.Contains(*g, upper) && strings.Contains(*g, lower) {
			*g = strings.ReplaceAll(*g, lower, upper)
		}
	}
	if w == lower {
		if !strings.Contains(*g, lower) && strings.Contains(*g, upper) {
			*g = strings.ReplaceAll(*g, upper, lower)
		}
	}

	// Does the value x repeat twice?
	ww := strings.Repeat(w, 2)
	if strings.Contains(*g, ww) {
		return myErrGO("The expression contains two literals together.")
	}

	// Remove all signs and numbers from a math expression.
	str, sSgn := *g, strings.Split("0123456789 ^+-."+w, "")
	for _, sgn := range sSgn {
		if strings.Contains(str, sgn) {
			str = strings.ReplaceAll(str, sgn, "")
		}
	}

	// Check for signs or literals other than normal.
	if str != "" {
		return myErrGO("The expression contains literals of const type.")
	}

	return nil
}

// Remove blank spaces between a sign and a number.
func directAllSign(g *string) {
	sSgn := []string{"+", "-"}
	for _, sgn := range sSgn {
		if strings.Contains(*g, sgn+" ") {
			*g = strings.ReplaceAll(*g, sgn+" ", sgn)
		}
	}
}

// Split str type math func into Kx^n expressions.
func splitY(g string) []string {
	// "+ n" => "+n"
	directAllSign(&g)

	// Divide the math func by its blanks.
	if strings.Contains(g, " ") {
		g = strings.Trim(g, " ")
		return strings.Split(g, " ")
	} else {
		return []string{g}
	}
}

// Finish building kxn for future operations.
func finishBuildingKxn(w string, kxn *string) {
	if !strings.Contains(*kxn, w+"^") {
		*kxn = strings.ReplaceAll(*kxn, w, w+"^1")

		if strings.Contains(*kxn, "-"+w) {
			*kxn = strings.ReplaceAll(*kxn, w, "1"+w)

		} else if strings.Contains(*kxn, "+"+w) {
			*kxn = strings.ReplaceAll(*kxn, w, "1"+w)
		}
	} else {
		if strings.Contains(*kxn, "-"+w) {
			*kxn = strings.ReplaceAll(*kxn, w, "1"+w)

		} else if strings.Contains(*kxn, "+"+w) {
			*kxn = strings.ReplaceAll(*kxn, w, "1"+w)
		}
	}
}

// kx^n → (kn)x^(n-1). For a single expression.
func theKxnRule(w string, kxn *string) {

	// All x or kx => kx^1
	finishBuildingKxn(w, kxn)

	// Index: "x" & Remove: "x^"
	iY := strings.Index(*kxn, w)
	*kxn = strings.ReplaceAll(*kxn, w+"^", "")
	sByte := []byte(*kxn)

	// ["k" => k] & ["n" => n]
	kByte, nByte := sByte[:iY], sByte[iY:]
	kStr, nStr := string(kByte), string(nByte)

	// Define a struct Kxn{}
	dy := dKxn{k: kStr, x: w, n: nStr}

	// Apply knxm method.
	*kxn = dy.knxm()

	// Set "w^0" to "".
	*kxn = strings.TrimSuffix(*kxn, w+"^0")
}

// Add blank spaces between a sign and a number.
func redirectASign(kxn string) string {
	// For all kxn other than "".
	if kxn != "" {

		// Add the missing signs.
		if strings.HasPrefix(kxn, "+") {
			kxn = " " + kxn
		} else if !strings.HasPrefix(kxn, "-") {
			kxn = " +" + kxn
		} else {
			kxn = " " + kxn
		}

		// Adjust signs  " + " or " - "
		sSign := []string{" +", " -"}
		for _, sign := range sSign {
			kxn = strings.Replace(kxn, sign, sign+" ", 1)
		}
	}

	return kxn
}

// kx^n → (kn)x^(n-1). For one or more expressions.
func DdKxn(w, g string) (string, string, error) {

	// Are there empty input fields?
	err := existsXY(w, g)
	if err != nil {
		return "", "", err
	}

	// Was the math func written correctly?
	err = correctlyXY(w, &g)
	if err != nil {
		return "", "", err
	}

	// Split on the expression Kx^n.
	sKxn := splitY(g)

	// Execute the rule Kx^n for each expression.
	// Doesn't take into account numbers without literal.
	var dg string
	for i, kxn := range sKxn {
		if strings.Contains(kxn, w) {

			// Add the + sign when it is only x
			moreLess := strings.HasPrefix(kxn, "+") || strings.HasPrefix(kxn, "-")
			if i == 0 && !moreLess {

				// Adjust the sign of the first kxn.
				sKxn[i] = "+" + kxn

				// Apply the rule.
				theKxnRule(w, &sKxn[i])

			} else {
				// Apply the rule.
				theKxnRule(w, &sKxn[i])
			}

			// Adjust signs.
			if i == 0 {
				dg += sKxn[i]
			} else {
				dg += redirectASign(sKxn[i])
			}

		} else {
			// Adjust the zero.
			if i == 0 {
				dg = "0"
			}
		}
	}

	// Sum all kxn with the same "x^n"
	err = AddKxn(&w, &dg)
	if err != nil {
		return "", "", err
	}

	return w, dg, nil
}
