package errorGenerator
func init() {
	//все соседние буквы на клавиатуре
	notCloseErrors = map[string]string{
		"й": "яыу",
		"ц": "фывк",
		"у": "ыва",
		"к": "вап",
		"е": "апр",
		"н": "про",
		"г": "рол",
		"ш": "олд",
		"щ": "лдж",
		"з": "джэ",
		"х": "жэ\\",
		"ъ": "жэ\\",

		"ф": "ця",
		"ы": "йцуячс",
		"в": "цукчсм",
		"а": "укесми",
		"п": "кенмит",
		"р": "енгить",
		"о": "нгштьб",
		"л": "гшщьбю",
		"д": "шщзбю.",
		"ж": "щзхбю.",
		"э": "зхъ.",

		"я": "ыв",
		"ч": "фыв",
		"с": "ыва",
		"м": "вап",
		"и": "апр",
		"т": "про",
		"ь": "рол",
		"б": "лдж",
		"ю": "джэ",
	}

	//только соседние буквы на одной строке
	closeErrors = map[string]string{
		"й": "цф", "q": "wa",
		"ц": "йу", "w": "qe",
		"у": "цк", "e": "wr",
		"к": "уе", "r": "et",
		"е": "кн", "t": "ry",
		"н": "ег", "y": "tu",
		"г": "нш", "u": "yi",
		"ш": "гщ", "i": "uo",
		"щ": "шз", "o": "ip",
		"з": "щх", "p": "o[",
		"х": "зъ",
		"ъ": "хэ",
	
		"ф": "ый", "a": "sq",
		"ы": "фв", "s": "ad",
		"в": "ыа", "d": "sf",
		"а": "вп", "f": "dg",
		"п": "ар", "g": "fh",
		"р": "по", "h": "gj",
		"о": "рл", "j": "hk",
		"л": "од", "k": "jl",
		"д": "лж", "l": "k;",
		"ж": "дэ",
		"э": "жх",
	
		"я": "чф", "z": "ax",
		"ч": "яс", "x": "zc",
		"с": "чм", "c": "xv",
		"м": "си", "v": "cb",
		"и": "мт", "b": "vn",
		"т": "иь", "n": "bm",
		"ь": "тб", "m": "n,",
		"б": "ью",
		"ю": "б.",
	}
}
