# Async sender.

## Описание проекта
Сервис позволяет осуществлять асинхнронную отправку запросов, не критичных 
к синхронности и не требующих получения ответа от системы.

Отлично подходит для отправки метрик c бекэнда в такие сервисы, как Google 
Analytics, Яндекс метрику или Mindbox. 

Со стороны сайта происходит локальный сетевой запрос (2-3мс) далее работа 
страницы не блокируется, продолжается выполнение.

## Запуск
1. Требуется IDE для разработки, разработка велась в https://www.jetbrains.com/go/
2. Компилируем бинарник в командной строке windows 
```bash
 set CGO_ENABLED=0; & set GOOS=linux; & go build -a -installsuffix cgo -o ./builds/asyncSender ./cmd/app/main.go
```
3. Либо запускаем в докере  
```bash
docker-compose up --build 
``` 
И вытаскиваем бинарник из контейнера, по пути /app/asyncSender.
4. Копируем файл [asyncsender.service_example](/asyncsender.service_example) в удобное место на сервере, 
c имененем asyncsender.service, указываем в нем нужные данные для запуска (путь до файла, имя пользователя, группу)
5. Создаем симлинк на этот файл в 
```bash
/etc/systemd/system/asyncsender.service
```
Просим админов разрешить start / stop своей службы для своего пользователя через sudo.

Подробная инструкция по [запуску сервисов через systemd](https://tuxotronic.org/post/go-service-over-systemd/)

## Использование

Пример запроса на PHP - cURL

```php
$data = ['test' => 'data'];
$curl = curl_init();

curl_setopt_array($curl, array(
    CURLOPT_URL => 'localhost:8065/api/send',
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_ENCODING => '',
    CURLOPT_MAXREDIRS => 10,
    CURLOPT_TIMEOUT => 0,
    CURLOPT_FOLLOWLOCATION => true,
    CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
    CURLOPT_CUSTOMREQUEST => 'POST',
    CURLOPT_POSTFIELDS => json_encode($data),
    CURLOPT_HTTPHEADER => [
        'Content-Type: application/json; charset=utf-8',
        'Accept: application/json',
        "Authorization: Mindbox secretKey='key'"
    ],
));

$response = curl_exec($curl);

curl_close($curl);
echo $response;
```

## Проверка работоспособности сервиса
Вы можете проверить работоспособность сервиса, отправив GET-запрос на `http://localhost:8065/api/health`. В ответе вернётся текущая дата, время и количество сообщений в очереди отправки.