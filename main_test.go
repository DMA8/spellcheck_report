package main

import (
	"fmt"
	"testing"

	"github.com/Saimunyz/speller" //спеллер
)
// go test -bench=. -benchmem -benchtime=100x main_test.go 
var speller1 *speller.Speller

func init() {
	speller1 = speller.NewSpeller("config.yaml")

	// load model
	err := speller1.LoadModel("models/model-without_singleWords.gz")
	if err != nil {
		fmt.Printf("No such file: %v\n", err)
		panic(err)
	}

}

//костюм гимнастический для девочки с бриджами
func BenchmarkSpellCheck(b *testing.B) {
	b.Run("test", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			speller1.SpellCorrect("кастюм")
		}
	})
}

func BenchmarkSpellCheck2(b *testing.B) {
	b.Run("test", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			speller1.SpellCorrect("сухрй шампкнь")
		}
	})
}

func BenchmarkSpellCheck3(b *testing.B) {
	b.Run("test", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			speller1.SpellCorrect("маслр для звгара")
		}
	})
}

func BenchmarkSpellCheck4(b *testing.B) {
	b.Run("test", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			speller1.SpellCorrect("зврядное устройствл для акумулятора")
		}
	})
}

func BenchmarkSpellCheck5(b *testing.B) {
	b.Run("test", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			speller1.SpellCorrect("ввто прлив для домпшних расиений")
		}
	})
}

func BenchmarkSpellCheck6(b *testing.B) {
	b.Run("test", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			speller1.SpellCorrect("сегнализация с автазапуском и обратнрй свясью")
		}
	})
}

func BenchmarkSpellCheck7(b *testing.B) {
	b.Run("test", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			speller1.SpellCorrect("ьраслет из натурпльного камьня и серибра пррбы")
		}
	})
}

func BenchmarkSpellCheck8(b *testing.B) {
	b.Run("test", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			speller1.SpellCorrect("мвска для кожы вакруг глас с липестками золотово")
		}
	})
}

func BenchmarkSpellCheck9(b *testing.B) {
	b.Run("test", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			speller1.SpellCorrect("мяхкая игрушька ростения из ростения протиф зомьби для дитей")
		}
	})
}

func BenchmarkSpellCheck10(b *testing.B) {
	b.Run("test", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			speller1.SpellCorrect("оверсайс фудболки висна летр для девачек подросков красная модноя молодешная")
		}
	})
}

func BenchmarkSpellCheck11(b *testing.B) {
	b.Run("test", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			speller1.SpellCorrect("польто жентское стеганое демисизонное с капушоном для прагулок с сабакой кросивое")
		}
	})
}

func BenchmarkSpellCheck12(b *testing.B) {
	b.Run("test", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			speller1.SpellCorrect("сандалии женскме нптуральная кожжа на пладформе италлия лето на липучьках на ноку")
		}
	})
}