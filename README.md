<h1>Ускоренная версия</h1>

С одной ошибкой 
|Speller           |Точность % (С токенайзером)| Логические ошибки | Ложные исправления | Однопоток Скорость запрос/секунда (запрос/миллисекунда)  | Многопоток Скорость запрос/секунда (запрос/миллисекунда)   | Память Мб 		 |  
|:------------------:|:------:|:----:|:---:|:----------------------------------------------------------:|:-----------------------------------------------------------:|:-------------------:|   
|Новый спеллер	(Обычная версия)	|	86.16% (92.23%) | 100%|	13% (Norm: 7%)| 161 запросов в секунду (0.16 запросов в миллисекунду) |  1820 запросов в секунду (1.8 запросов в миллисекунду) |2824 Мб (5206 мб в многопотоке)	|


--------------------------------------------------------

С двумя ошибкой 
|Speller           |Точность % (С токенайзером)| Логические ошибки | Ложные исправления | Однопоток Скорость запрос/секунда (запрос/миллисекунда)  | Многопоток Скорость запрос/секунда (запрос/миллисекунда)   | Память Мб 		 |  
|:------------------:|:------:|:----:|:---:|:----------------------------------------------------------:|:-----------------------------------------------------------:|:-------------------:|   
|Новый спеллер	(Обычная версия)	|	79.79% (83.89%) | 100%|	13% (Norm: 7%)| 183 запросов в секунду (0.18 запросов в миллисекунду) |  2455 запросов в секунду (2.4 запросов в миллисекунду) |2823 Мб (4910 мб в многопотоке)	|

--------------------------------------------------------

<h2> Альтернативная версия </h2>

Пропускаем слова, которые есть в словаре частот, чтобы не тратить время на их исправления.<br>
Из-за этого теряется не отрабатывает контекст в запросах типа "подушка розовая" "подушка разовая". Однако подобный подход ускоряет работу спеллера в запросах с небольшим числом ошибок (слов, которых нет в словаре частот).<br>
В таблице выше, эффективность работы спеллера в сложных кейсах, когда слова написаны правильно, но не подходят по смыслу отображена в колонке Logic. Данная версия спеллера справляется только в 20% кейсов.<br>
AllWords/ErrorWords - отношение всех слов запроса к ошибочным (с округлением к большему), но не меньше 1 ошибочного слова на запрос. Например, если слов в запросе 5, а соотношение равно 2м, то будет сгенерировано 3 ошибки в случайных позициях (псевдослучайных)<br>


| AllWords/ErrorWords | Rate % | NormRate % | RPS  | errorsInWord |
|---------------------|--------|------------|------|--------------|
| 1 (каждое слово запроса с ошибкой)                  | 82.69  | 87.5       | 4241 | 1            |
| 2                   | 86.37  | 90.72      | 4594 | 1            |
| 3                   | 88.89  | 92.49      | 4914 | 1            |
| 4                   | 89.91  | 93.3       | 5352 | 1            |
| 5                   | 90.87  | 93.35      | 5650 | 1            |
| 1 (каждое слово запроса с ошибкой)                    | 68.21  | 71.85      | 5459 | 2            |
| 2                   | 78.79  | 82.20      | 4824 | 2            |
| 3                   | 82.41  | 85.57      | 5221 | 2            |
| 4                   | 83.60  | 86.75      | 5716 | 2            |
| 5                   | 84.19  | 87.32      | 5945 | 2            |
| Запросы без ошибок  | 94.59  | 95.81      | 126867 | -------      |

<h1>Результаты тестов (50к тестовых запросов old)</h1>
С одной ошибкой  

|Speller           |Точность % (С лемматизацией)| Однопоток Скорость запрос/секунда (запрос/миллисекунда)  | Многопоток Скорость запрос/секунда (запрос/миллисекунда)   | Память Мб 		 |  
|:------------------:|:---------:|:----------------------------------------------------------:|:-----------------------------------------------------------:|:-------------------:|   
|Новый спеллер		|	89.58% (92.74%)	| 75 запросов в секунду (0.07 запросов в миллисекунду) | 356 запросов в секунду (0.35 запросов в миллисекунду) |3029 Мб (5858 мб в многопотоке)	|
|Yandex     		|	78.30% (80.40%)	|   --   | 68 запросов в секунду (0.06 запросов в миллисекунду)	| --   |
|Текущий спеллер	|8.71%| 125445 запросов в секунду ( 125 запросов в миллисекунду)	  | --	|798 Мб|

--------------------------------------------------------

С двумя ошибками
|Speller           |Точность % (С лемматизацией)| Однопоток Скорость запрос/секунда (запрос/миллисекунда)  | Многопоток Скорость запрос/секунда (запрос/миллисекунда)   | Память Мб 		 |  
|:------------------:|:---------:|:----------------------------------------------------------:|:-----------------------------------------------------------:|:-------------------:|   
|Новый спеллер		|	87.35% (90.52%)	| 171 запросов в секунду (0.17 запросов в миллисекунду)     |	775 запросов в секунду (0.7 запросов в миллисекунду)     |	 3010 Мб  (5596 мб в многопотоке)    |
|Yandex     		|	13.49% (14.34%)	|	--  |  77 запросов в секунду (0.07 запросов в миллисекунду)    	|	--	  |
|Текущий спеллер	|0.08%|96721 запросов в секунду ( 96 запросов в миллисекунду)	  |	--	|	795 Мб|

 Производительность ЯндексСпеллера тестилась 12 воркерами, потому что ждать завершения синхронных 50к тестов очень долго.<br>
 Производительность нашего спеллера в многопотоке тестилась с 12 воркерами. Увеличение кол-ва воркеров не дает ощутимого преимущества.<br>
 Продовский спеллер такой быстрый, потому что он представляет собой 2 хэштаблицы.<br><br>
С лемматизацией - процент успешных исправлений спеллера, при сравнении лемм предложенного варианта спеллера и эталонного запроса. Лемматизация после получения исправления для приведения к одной форме для тестов.  
Ошибки генерируются псевдослучайно. Рандомная буква слова заменяется на соседнюю букву на клавиатуре. На кажlдый тестовый поисковый запрос создается 5 ошибочных вариантов<br>
Если длина слова 3-4 буквы, то максимум генерируется 1 ошибка.<br>
<hr>
<h3>Бенчмарки</h3>
Результаты бенчмарков в папке bench/<br>
Было проведено 12 тестов с разными размерами запросов (от 1 до 12)<br>
Продовский спеллер гораздо быстрее, экономичнее по памяти, но его успешность исправлений ~1%.
<hr>
<h3>Комментарии</h3>
<p><i>Токенайзер</i></p>
Токенайзер добавить не удалось, точнее с ним тестировали он все нормализует, но таким образом ломается контекст, некоторые слова становятся частотнее других, таким образом все только портится, поэтому пока что токенайзер для обучения нам не подходит. Можно его попробовать  использовать потом для определения частей речи и как-то это дальше добавлять в спеллер, просто как идея, пока не пробовали и не обдумывали до конца.<br>
Использовали токенайзер для тестов, чтобы убрать проблемы в разных формах слов. Точно можно сказать, что для улучшения нужно собирать и предобрабатывать тренируемый текст, так как после смены его на хитовые запросы все стало в разы лучше, чем то что мы собирали с карточек.<br><br>
<p><i>Проблема длинных запросов</i></p>
когда получаем запрос из 7+ слов то спеллер генерит кучу перестановок с возможными ошибками для каждого слова, тем самым убивается и время и память очень сильно. Поэтому, пока делим подобные запросы на подзапросы поменьше, по 3 слова, из-за этого слегка может тоже страдать контекст, так как не видит полной картины, как это сделать более правильно пока не придумали.<br><br>
<p><i>Оптимизации памяти</i></p>
Когда обучали на карточках, модель занимала около 4-5 ГБ, там было около 10 млн строк. Теперь после того как поменяли построение n-грамм по строкам, чтобы не создавались лишние и не используемые, память слегка сократилась (Изначально спеллер воспринимал весь учебный текст - как одну большую строку, тем самым генерив огромное число нерелевантные n-граммы). С новой моделью из хитовых запросов, она 1.2 млн строк, все стало весить в памяти около 1-2 ГБ.<br>
Так же изменили обработку коротких слов (до 3х включительно букв). Короткие слова не добавляются в дерево ngramm (не участвуют в определении контекста) и исправляются изолированно. Это позволило облегчить модель и ускорило обработку длинных запросов с короткими словами с воставе.
<hr>

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
