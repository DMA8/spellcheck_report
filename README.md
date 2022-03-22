<p>1) Без токенизации учебной модели и ожидаемого ответа : </p>
	Results: <br>
	TotalTests: 7388<br>
	SpellerRate 50.05%, YandexRate 9.77% <br><br>
Logs: <br>
Все кейсы: ./logs/not_normalized_logs.txt <br>
Яндекс справился, а наш спеллер нет: ./logs/not_normalized_yaRigth_spellWrong_log.txt

------------------------------------------------------------------------------------------

<p>2) Токенизация учебной модели, но без токенизации ожидаемого ответа : </p>
		Results: <br>
		TotalTests: 7537 <br>
		SpellerRate 19.56%, YandexRate 9.58% <br><br>
Все кейсы: ./logs/normilizedModel_notNormalizedTestWord_logs.txt<br>
Яндекс справился, а наш спеллер нет: ./logs/normilizedModel_notNormalizedTestWord_yaRigth_spellWrong_log.txt

------------------------------------------------------------------------------------------
<p>3) Токенизация учебной модели и токенизация ожидаемого ответа: </p>
		Results:<br>
		TotalTests: 7429<br>
		SpellerRate 41.59%, YandexRate 6.97%<br><br>
Все кейсы: ./logs/normalizedModel_normalizedTestWord_logs.txt<br>
Яндекс справился, а наш спеллер нет: ./logs/normalizedModel_normalizedTestWord_yaRigth_spellWrong_log.txt

------------------------------------------------------------------------------------------
