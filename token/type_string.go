// Code generated by "stringer -type Type"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EOF-0]
	_ = x[LeftParen-1]
	_ = x[RightParen-2]
	_ = x[LeftBrack-3]
	_ = x[RightBrack-4]
	_ = x[LeftBrace-5]
	_ = x[RightBrace-6]
	_ = x[Comma-7]
	_ = x[Dot-8]
	_ = x[Minus-9]
	_ = x[Plus-10]
	_ = x[Star-11]
	_ = x[Slash-12]
	_ = x[Semicolon-13]
	_ = x[Bang-14]
	_ = x[BangEqual-15]
	_ = x[Equal-16]
	_ = x[EqualEqual-17]
	_ = x[Less-18]
	_ = x[LessEqual-19]
	_ = x[Greater-20]
	_ = x[GreaterEqual-21]
	_ = x[And-22]
	_ = x[Class-23]
	_ = x[Else-24]
	_ = x[False-25]
	_ = x[For-26]
	_ = x[Fun-27]
	_ = x[If-28]
	_ = x[Nil-29]
	_ = x[Or-30]
	_ = x[Print-31]
	_ = x[Return-32]
	_ = x[Super-33]
	_ = x[This-34]
	_ = x[True-35]
	_ = x[Var-36]
	_ = x[While-37]
	_ = x[String-38]
	_ = x[InvalidString-39]
	_ = x[Number-40]
	_ = x[InvalidNumber-41]
	_ = x[Identifier-42]
	_ = x[Other-43]
}

const _Type_name = "EOFLeftParenRightParenLeftBrackRightBrackLeftBraceRightBraceCommaDotMinusPlusStarSlashSemicolonBangBangEqualEqualEqualEqualLessLessEqualGreaterGreaterEqualAndClassElseFalseForFunIfNilOrPrintReturnSuperThisTrueVarWhileStringInvalidStringNumberInvalidNumberIdentifierOther"

var _Type_index = [...]uint16{0, 3, 12, 22, 31, 41, 50, 60, 65, 68, 73, 77, 81, 86, 95, 99, 108, 113, 123, 127, 136, 143, 155, 158, 163, 167, 172, 175, 178, 180, 183, 185, 190, 196, 201, 205, 209, 212, 217, 223, 236, 242, 255, 265, 270}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}