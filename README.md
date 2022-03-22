Без токенизации учебной модели и ожидаемого ответа :
		Results:
		 TotalTests: 7388
		 SpellerRate 50.05%, YandexRate 9.77%
Все кейсы: ./logs/not_normalized_logs.txt
Яндекс справился, а наш спеллер нет: ./logs/not_normalized_yaRigth_spellWrong_log.txt


Токенизация учебной модели, но без токенизации ожидаемого ответа :
		Results:
		 TotalTests: 7537
		 SpellerRate 19.56%, YandexRate 9.58%
Все кейсы: ./logs/normilizedModel_notNormalizedTestWord_logs.txt
Яндекс справился, а наш спеллер нет: ./logs/normilizedModel_notNormalizedTestWord_yaRigth_spellWrong_log.txt

Токенизация учебной модели и токенизация ожидаемого ответа:
		Results:
		 TotalTests: 7429
		 SpellerRate 41.59%, YandexRate 6.97%
Все кейсы: ./logs/normalizedModel_normalizedTestWord_logs.txt
Яндекс справился, а наш спеллер нет: ./logs/normalizedModel_normalizedTestWord_yaRigth_spellWrong_log.txt