<p>1) Без токенизации учебной модели и ожидаемого ответа : </p>
	Results: <br>
	TotalTests: 7388<br>
	SpellerRate 50.05%, YandexRate 9.77% <br><br>
Logs: <br>
Все кейсы: ./logs/not_normalized_logs.txt <br>
Яндекс справился, а наш спеллер нет: ./logs/not_normalized_yaRigth_spellWrong_log.txt

------------------------------------------------------------------------------------------

<p>2) Токенизация учебной модели, но без токенизации ожидаемого ответа : </p>
		<p>Results: </p>
		<p>TotalTests: 7537 </p>
		<p>SpellerRate 19.56%, YandexRate 9.58% </p>
Все кейсы: ./logs/normilizedModel_notNormalizedTestWord_logs.txt
Яндекс справился, а наш спеллер нет: ./logs/normilizedModel_notNormalizedTestWord_yaRigth_spellWrong_log.txt

------------------------------------------------------------------------------------------
<p>3) Токенизация учебной модели и токенизация ожидаемого ответа: </p>
		<p>Results:
		<p>TotalTests: 7429
		<p>SpellerRate 41.59%, YandexRate 6.97%
Все кейсы: ./logs/normalizedModel_normalizedTestWord_logs.txt
Яндекс справился, а наш спеллер нет: ./logs/normalizedModel_normalizedTestWord_yaRigth_spellWrong_log.txt

------------------------------------------------------------------------------------------
