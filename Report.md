# Домашнее задание №10

Начальный прогон бенчмарков.

  ![0_init](img/1.png)

1. Оптимизация вызова функции regexp.MatchString
  - При отображении функций отсортированных по потреблению cpu видно, что функции для работы с регулярными выражениями занимают больше всего процессорного времени (MatchString и Compile).

  ![1_regex_pprof](img/2.png)

  - Также при исследовании функции FastSearch построчно видно, что на строчках 60 и 82, где непосредственно вызывается MatchString, потрачено много времени.    

  ![1_regex_list](img/3.png)

  - Главная причина такого потребления заключается в том, что функция regexp.MatchString вызывается в цикле большое количество раз, и под капотом включает в себя компилирование регулярного выражения с заданным паттерном. Компилирование регулярки очень затратная операция, да и сам матчинг тоже не быстрый. Данный момент можно было бы оптимизировать путем прекомпилирования данных регулярок, т.е. сделать их глобальными переменными. Однако данные регулярные выражения очень просты, и целесообразнее будет вместо них применить strings.Contains. Кроме того, regexp.MatchString потенциально может быть источником ошибки.

  - После первой же оптимизации виден прирост в производительности.

  ![1_regex_result](img/4.png)

2. Оптимизация функции json.Unmarshal
  - Проделав операции, аналогичные тем, что были в первом пункте, видим, что теперь наиболее затратная операция - json.Unmarshal.
  
  ![2_unmarshal_pprof](img/5.png)

  - Видно, что json.Unmarshal вызывается на строчке 36, и затрачивает огромное количество ресурсов

  ![2_unmarshal_list](img/6.png)

  - Причина большого потребления ресурсов данной функции в том, что она внутри использует рефлексию. Поэтому применим кодогенератор easyjson, который нам сгенерирует высокопроизводительный и явный код для десериализации юзеров в структуру. Прирост производительности также заметен.
  
  ![2_unmarshal_result](img/7.png)

3. Оптимизация исключением лишних утверждений типов (type assertions) и заменой способа хранения юзеров с мапы(map[string]interface{}) на структуру юзера.
  - Данный пункт является следствием предыдущего, поскольку применив кодогенерацию, потребность в данных вещах отпала, и это также благоприятно сказывается на потреблении cpu и памяти.

4. Оптимизация вызова функции ReplaceAllString
  - Данная оптимизация проводится аналогичным образом, как в первом пункте. Причины все те же, regexp можно заменить обычной функцией из пакета strings.

  ![4_regex_list](img/8.png)

  - После оптимизации видим небольшие улучшения

  ![4_regex_result](img/9.png)

5. Оптимизация вызова outil.ReadAll(file)
  - При исследовании профиля памяти сразу бросается в глаза чрезмерное потребление памяти функцией outil.ReadAll(file). Проблема заключается в том, что нет нужды считывать весь файл в память, поскольку работа идет лишь с одной строкой(пользователем). При увеличении размера входных данных данная проблема усугубится еще сильнее.

  ![5_readall_pprop](img/10.png)
  
  - Исправляем данный момент путем чтения за раз одной строки и последующей ее обработки. В дополнении, данная оптимизация убирает нужду в функции strings.Split, а также преобразование байтов к строке внутри нее, которые нагружали cpu и особенно память. Также устраняется преобразование строки в слайс байтов при вызове Unmarshal и перевыделение памяти слайса users при добавлении в него юзеров, поскольку слайс при создании не был преаллоцирован на нужный размер.
  - После данной оптимизации видим серьезное улучшение всех показателей.

  ![5_readall_result](img/11.png)

6. Оптимизация конкатенации строк
  - Опять исследуем профиль памяти и видим, что на строчках 91 и 94 очень большое потребление памяти. Все дело в конкатенации строк. Данная операция очень ресурсоемкая поскольку каждый раз при склеивании строк происходит выделение памяти для новой строки и затем копирование данных соединяемых строк. С увеличением строки это бьет все больнее и больнее. Кроме того излишне нагружаем GC.

  ![6_list_strconcat](img/12.png)

  - Для оптимизации этого момента применим strings.Builder. А вместо конкатенации в функции fmt.Fprintln используем fmt.Fprintf и распечатаем строку через %s
  - Видим улучшения по памяти
  
  ![6_strconcat_result](img/13.png)

7. Оптимизация аллокации структуры юзера
  - Поскольку в один момент мы обрабатываем одну строку/юзера, то нет смысла выдялять память под юзера на каждой итерации цикла. Достаточно создать переменную вне цикла и занулять значение в конце цикла.
 
  ![7_list_user](img/14.png)
  
  Результат оптимизации заметен

  ![7_user_result](img/15.png)

8. Оптимизация поиска встречавшихся ранее браузеров и их добавление; удаление двойного итерирования по слайсу с браузерами
  - При исследовании профиля cpu видны затраты при итерировании по слайсу со встретившимися ранее браузерами. Также данный слайс хранит уникальные значения, и для чтобы проверить встречался ли ранее браузер необходимо каждый раз итерироваться по всему слайсу. Это линейная зависимость, а можно сделать константую, использовав мапу для этих целей. Также заранее преалоцируем некоторое количество памяти. Также несколько зарефакторим данные участки кода и удалим ненужные переменные uniqueBrowsers и notSeenBefore. Заодно уберем двойное прохождение по слайсу browsers.
  
  ![8_list_browsers](img/16.png)

  - В результате немного улучшилась производительность
  
  ![8_browsers_result](img/17.png)

### Сравнение с BenchmarkSolution
|                     |     |               |             |
|---------------------|-----|---------------|-------------|
| BenchmarkSolution-8 | 500 | 2782432 ns/op | 559910 B/op | 10422 allocs/op   
| BenchmarkFast-4     | 658 | 1728695 ns/op | 496272 B/op | 6478 allocs/op