package main

import (
    "fmt"
    "log"
    "github.com/streadway/amqp"
)

// Функция для обработки ошибок
func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}

func main() {
    // Подключение к RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:15672/")
    failOnError(err, "Ошибка подключения к брокеру")
    defer conn.Close() // Закрытие соединения при выходе из функции

    ch, err := conn.Channel()
    failOnError(err, "Ошибка открытия канала")
    defer ch.Close() // Закрытие канала при выходе из функции

    // Объявление обменника
    err = ch.ExchangeDeclare(
        "orders",   // имя
        "topic",    // тип
        true,       // durable
        false,      // auto-deleted
        false,      // internal
        false,      // no-wait
        nil,        // аргументы
    )
    failOnError(err, "Ошибка создания обменника")

    // Объявление очереди
    q, err := ch.QueueDeclare(
        "order_queue", // имя
        true,          // durable
        false,         // auto-delete
        false,         // exclusive
        false,         // no-wait
        nil,           // аргументы
    )
    failOnError(err, "Ошибка создания очереди")

    // Привязка очереди к обменнику
    err = ch.QueueBind(
        q.Name,       // имя очереди
        "order.*",    // ключ маршрутизации
        "orders",     // имя обменника
        false,
        nil,
    )
    failOnError(err, "Ошибка привязки очереди")

    // Публикация сообщения
    order := []byte(`{"id": "1", "item": "Pizza", "quantity": 2}`)
    err = ch.Publish(
        "orders",         // обменник
        "order.create",   // ключ маршрутизации
        false,            // обязательный
        false,            // немедленный
        amqp.Publishing{
            ContentType: "application/json",
            Body:        order,
        },
    )
    failOnError(err, "Ошибка публикации сообщения")
    fmt.Println("Сообщение отправлено:", string(order))

    // Ожидание сообщений
    msgs, err := ch.Consume(
        q.Name, // имя очереди
        "",     // consumer
        true,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // аргументы
    )
    failOnError(err, "Ошибка подписки на очередь")

    // Обработка сообщений в отдельной горутине
    go func() {
        for d := range msgs {
            fmt.Printf("Получено сообщение: %s\n", d.Body)
            // Здесь можно добавить логику обработки заказа
        }
    }()

    fmt.Println("Ожидание сообщений...")
    select {} // Блокировка основного потока, чтобы программа не завершилась
}