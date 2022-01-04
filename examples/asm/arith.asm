. program
arith	START	0

. vsota
		LDA		x
		ADD		y
		STA		sum

. razlika
		LDA		x
		SUB		y
		STA		diff

. zmnozek
		LDA		x
		MUL		y
		STA		prod

. kvocient
		LDA		x
		DIV		y
		STA		quot

. modulo
modulo	LDA		x
		COMP	y
		JGT		raz
		STA		mod
		J		halt

raz		SUB		y
		STA		x
		J		modulo

halt	J		halt

. podatki
x		WORD	10
y		WORD	3
sum		RESW	1
diff	RESW	1
prod	RESW	1
quot	RESW	1
mod		RESW	1
