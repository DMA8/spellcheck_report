<h1>Результаты тестов</h1>
Две ошибки в тестовом запросе на каждое слово <br>
TotalTests: 3005<br>
Полностью исправленные запросы<br>
SpellerRateQuery 66.99% (Norm: 71.11%),  YandexRateQuery 28.05% (Norm: 29.45%)<br>
Исправленно слов из запросов:<br>
Total words: 9160, SpellerRateWords 86.35%, YandexRateWords 54.57%<br>
<br>
Одна ошибка в тестовом запросе на каждое слово<br>
TotalTests: 3005<br>
Полностью исправленные запросы:<br>
SpellerRateQuery 63.03% (Norm: 70.42%),  YandexRateQuery 73.21% (Norm: 75.04%)<br>
Исправленно слов из запросов:<br>
Total words: 9240, SpellerRateWords 83.00%, YandexRateWords 88.47%<br><br>
Ошибки генерируются псевдослучайно. Рандомная буква слова заменяется на соседнюю букву на клавиатуре. На кадждый поисковый запрос создается 5 ошибочных вариантов<br>
Если длина слова 3-4 буквы, то максимум генерируется 1 ошибка.<br>
<hr>


<h3>Комментарии</h3>
<p><i>Токенайзер</i></p>
Токенайзер добавить не удалось, точнее с ним тестировали он все нормализует, но таким образом ломается контекст, некоторые слова становятся частотнее других, таким образом все только портится, поэтому пока что токенайзер для обучения нам не подходит. Можно его попробовать  использовать потом для определения частей речи и как-то это дальше добавлять в спеллер, просто как идея, пока не пробовали и не обдумывали до конца.<br>
Использовали токенайзер для тестов, чтобы убрать проблемы в разных формах слов. Точно можно сказать, что для улучшения нужно собирать и предобрабатывать тренируемый текст, так как после смены его на хитовые запросы все стало в разы лучше, чем то что мы собирали с карточек.<br><br>
<p><i>Проблема длинных запросов</i></p>
когда получаем запрос из 7+ слов то спеллер генерит кучу перестановок с возможными ошибками для каждого слова, тем самым убивается и время и память очень сильно. Поэтому, пока делим подобные запросы на подзапросы поменьше, по 3 слова, из-за этого слегка может тоже страдать контекст, так как не видит полной картины, как это сделать более правильно пока не придумали.<br><br>
<p><i>Оптимизации памяти</i></p>
Когда обучали на карточках, модель занимала около 4-5 ГБ, там было около 10 млн строк. Теперь после того как поменяли построение n-грамм по строкам, чтобы не создавались лишние и не используемые, память слегка сократилась (Изначально спеллер воспринимал весь учебный текст - как одну большую строку, тем самым генерив огромное число нерелевантные n-граммы). С новой моделью из хитовых запросов, она 1.2 млн строк, все стало весить в памяти около 1-2 ГБ.<hr>
<h1>Интерпретация логов</h1>
<br>
Логи тестирования спеллера приведены в папках oneErrorLogs/ и twoErrorsLogs/.<br>
<h3>1. oneError.txt (или twoError.txt)</h3>
содержит все проведенные тесткейсы. Тесткейсы разделены между собой по тестируемому слову, которое указывается на первой строке блока. На последующих 5 строках блока указывается:<br> 
"generated error is:" - запрос с одной или двумя искусственно сгенерированных ошибки на каждое слово<br>
"S:" -  вариант исправления ошибки, предложенный нашим спеллером<br>
"Y:" - вариант исправления ошибки, предложенный Яндекс спеллером<br>
На последней строке блока - число успешных исправлений спеллеров в контексте текущего блока<br>
<h3>2. yaRight_spellWrong_log.txt</h3>
содержит тесткейсы, в которых ЯндексСпеллер исправляет корректно, а наш спеллер - нет.<br> В первой строке блока: <br>
(средство для чистки кроссовок -> спедмтво ддя чмстуи кроссовре) - (`исходный запрос -> сгенерированный ошибочный запрос`) <br>
"yaSuceed:" (вариант, предложенный Яндекс спеллером)<br>
"spellerFail:" - вариант, предложенный нашим спеллером<br>
на последующих строках детально разбираются ошибки спеллера. Последующие строки состоят из:<br>
"Error:" - сгенерированная ошибка, с которой не справлися спеллер<br>
"Expected:" - исходное слово, из которого сгенерировали ошибку. freq - частота слова в учебной модели. diffRunes - сколько букв отличаются у Error и Expected<br>
"SpellerSuggest:" - неправильный вариант исправления, предложенный нашим спеллером (freq и diffRunes для слова, предложенного спеллером).<br>
<i>При разборе логов, обращая внимание на freq и diffRunes, можно понять, почему спеллер принял решение в пользу того, или иного слова. Нельзя ожидать от спеллера, что он справится с исправлениями, если у слова Expected freq равно 0.</i>
<h4>3. yaRight_spellWrong_log.txt</h4>
Содержит тесткейсы, в которых спеллер справляется, а ЯндексСпеллер - нет<br>
"W:" - исходное слово; "E:" - сгенерированная ошибка, "Y:" - ошибочный вариант ЯндексСпеллера, "S:" - правильный вариант спеллера<br>
<h4>4.bothWrongLog.txt</h4>
Тесткейсы, в которых оба спеллера не справились<br>
"Expected:" - исходное слово; "Error" - сгенерированная ошибка, "SpellerSuggest:" - ошибочный вариант спеллера "YandexSuggest:" - ошибочный вариант ЯндексСпеллера<br>
<h5>5.notmalize....txt</h5>
Файлы, начинающиеся с "normalize" попадают тесткейсы, в которых спеллеры не справляются. Слова, которые они не смогли исправить пропускаются через токенайзер и преобразуются в леммы(начальная форма слова). Затем полученные леммы сравниваются с леммами от исходного слова. Это позволяет уточнить эффективность спеллеров, так как мы не считаем неверную форму одного и того же слова за ошибку.
