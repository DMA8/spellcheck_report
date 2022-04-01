<h1>Результаты тестов</h1>
по 2 ошибки в тестовом запросе на каждое слово (если длина слова 3-4 буквы, то генерируется 1 ошибка)<hr>
Results:<br>
Полностью исправленные запросы(TotalTests: 3005):<br>
 SpellerRate 66.99% (Norm: 71.11%),  YandexRate 28.05% (Norm: 29.45%)<br>
Исправленно слов из запросов:<br>
Total words: 9160, SpellerRate 86.35%, YandexRate 54.57%<br>
<br>
по 1 ошибке в тестовом запросе на каждое слово (если длина слова меньше 3 букв, то ошибка не генерируется)<br>
Results:<br>
Полностью исправленные запросы(TotalTests: 3005):<br>
SpellerRate 63.03% (Norm: 70.42%),  YandexRate 73.21% (Norm: 75.04%)<br>
Исправленно слов из запросов:<br>
Total words: 9240, SpellerRate 83.00%, YandexRate 88.47%<br>
<h1>Интерпретация логов</h1>
Ошибки генерируются псевдослучайно. Рандомная буква слова заменяется на соседнюю букву на клавиатуре. На кадждый поисковый запрос создается 5 ошибочных вариантов<br><br>
Логи тестирования спеллера приведены в папках **oneErrorLogs/** и **twoErrorsLogs/**.<br><br>
**oneError.txt** (или **twoError.txt**) - содержит все проведенные тесткейсы. Тесткейсы разделены между собой по тестируемому слову, которое указывается на первой строке блока. На последующих 5 строках блока указывается:<br> 
"***generated error is:***" - запрос с одной или двумя искусственно сгенерированных ошибки на каждое слово<br>
"***S:***" -  вариант исправления ошибки, предложенный нашим спеллером<br>
"***Y:***" - вариант исправления ошибки, предложенный Яндекс спеллером<br>
На последней строке блока - число успешных исправлений спеллеров в контексте текущего блока<br><br>
