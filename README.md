ТЕСТ БРЕНДОВ:

<p>Бренды на английском</p>
Results:<br>
 TotalTests: 8005<br>
 SpellerRate 43.44%, YandexRate 10.93%<br>
 
 YandexFails 7130 SpellerRight 2854 SpellerRate 40.03% (Спеллер исправил 40% ошибок яндекса)<br>
Logs: <br>
Все кейсы: ./logs/EngBrandsLogs.txt <br>
Яндекс справился, а наш спеллер нет: ./logs/EngBrandsYaRigth_spellWrong_log.txt<br>
Яндекс ошибся, а наш спеллер смог исправить: ./logs/EngBrandsSpellerRightYandexWrong.txt<br>

----------------------------------------------------------------------------------------

<p>Бренды на русском</p>
Results:<br>
 TotalTests: 8005<br>
 SpellerRate 59.14%, YandexRate 28.31%<br>
 YandexFails 5739 SpellerRight 2816 SpellerRate 49.07%  (Спеллер исправил 49% ошибок яндекса)<br>
<br>
Logs: <br>
Все кейсы: ./logs/RU_Brand_Logs.txt <br>
Яндекс справился, а наш спеллер нет: ./logs/RU_BrandLogs_YaRigth_spellWrong_log.txt<br>
Яндекс ошибся, а наш спеллер смог исправить: ./logs/RU_Brand_LogsSpellerRightYandexWrong.txt<br>

==========================================================================================
ТЕСТ ПОИСКОВЫХ ЗАПРОСОВ ИЗ НАЗВАНИЙ КАРТОЧЕК:

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
