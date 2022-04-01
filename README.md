<h1>Результаты тестов</h1>
Две ошибки в тестовом запросе на каждое слово (если длина слова 3-4 буквы, то генерируется 1 ошибка)<br>
TotalTests: 3005<br>
Полностью исправленные запросы<br>
SpellerRate 66.99% (Norm: 71.11%),  YandexRate 28.05% (Norm: 29.45%)<br>
Исправленно слов из запросов:<br>
Total words: 9160, SpellerRate 86.35%, YandexRate 54.57%<br>
<br>
Одна ошибка в тестовом запросе на каждое слово (если длина слова меньше 3 букв, то ошибка не генерируется)<br>
TotalTests: 3005<br>
Полностью исправленные запросы:<br>
SpellerRate 63.03% (Norm: 70.42%),  YandexRate 73.21% (Norm: 75.04%)<br>
Исправленно слов из запросов:<br>
Total words: 9240, SpellerRate 83.00%, YandexRate 88.47%<br>
<h1>Интерпретация логов</h1>
Ошибки генерируются псевдослучайно. Рандомная буква слова заменяется на соседнюю букву на клавиатуре. На кадждый поисковый запрос создается 5 ошибочных вариантов<br><br>
Логи тестирования спеллера приведены в папках oneErrorLogs/ и twoErrorsLogs/.<br><br>
<h3>1. oneError.txt (или twoError.txt)</h3>
содержит все проведенные тесткейсы. Тесткейсы разделены между собой по тестируемому слову, которое указывается на первой строке блока. На последующих 5 строках блока указывается:<br> 
"generated error is:" - запрос с одной или двумя искусственно сгенерированных ошибки на каждое слово<br>
"S:" -  вариант исправления ошибки, предложенный нашим спеллером<br>
"Y:" - вариант исправления ошибки, предложенный Яндекс спеллером<br>
На последней строке блока - число успешных исправлений спеллеров в контексте текущего блока<br><br>
<h3>2. yaRight_spellWrong_log.txt</h3>
содержит тесткейсы, в которых ЯндексСпеллер исправляет корректно, а наш спеллер - нет.<br> В первой строке блока: <br>
(средство для чистки кроссовок -> спедмтво ддя чмстуи кроссовре) - (`исходный запрос -> сгенерированный ошибочный запрос`) <br>
"yaSuceed:" (вариант, предложенный Яндекс спеллером)<br>
"spellerFail:" - вариант, предложенный нашим спеллером<br>
на последующих строках детально разбираются ошибки спеллера. Последующие строки состоят из:<br>
"Error:" - сгенерированная ошибка, с которой не справлися спеллер<br>
"Expected:" - исходное слово, из которого сгенерировали ошибку. freq - частота слова в учебной модели. diffRunes - сколько букв отличаются у Error и Expected<br>
"SpellerSuggest:" - неправильный вариант исправления, предложенный нашим спеллером (freq и diffRunes для слова, предложенного спеллером).<br>
При разборе логов, обращая внимание на freq и diffRunes, можно понять, почему спеллер принял решение в пользу того, или иного слова. Нельзя ожидать от спеллера, что он справится с исправлениями, если у слова Expected freq равно 0.
<h4>3. yaRight_spellWrong_log.txt</h4>
Содержит тесткейсы, в которых спеллер справляется, а ЯндексСпеллер - нет<br>
"W:" - исходное слово; "E:" - сгенерированная ошибка, "Y:" - ошибочный вариант ЯндексСпеллера, "S:" - правильный вариант спеллера<br>
