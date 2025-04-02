package main

import (
    "log"
    "net/smtp"
    "github.com/streadway/amqp"
    "encoding/json"
)

type Order struct {
    ID       string `json:"id"`
    Item     string `json:"item"`
    Quantity int    `json:"quantity"`
}

func sendEmail(order Order) {
    from := "chipavoyvoy@gmail.com"
    password := "Rolan.2005"
    to := "alisertagirovich52@gmail.com"
    subject := "New Order Received"
    body := "Order ID: " + order.ID + "\nItem: " + order.Item + "\nQuantity: " + string(order.Quantity)

    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    err := smtp.SendMail("smtp.example.com:587",
        smtp.PlainAuth("", from, password, "smtp.example.com"),
        from, []string{to}, []byte(msg))

    if err != nil {
        log.Printf("Ошибка отправки почты: %s", err)
    } else {
        log.Println("Письмо отправлено успешно")
    }
}

func main() {
    // Подключение к RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:15672/")
    if err != nil {
        log.Fatalf("Ошибка подключения к брокеру: %s", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Ошибка открытия канала: %s", err)
    }
    defer ch.Close()

    // Объявление очереди
    q, err := ch.QueueDeclare(
        "order_queue", // имя
        true,          // durable
        false,         // auto-delete
        false,         // exclusive
        false,         // no-wait
        nil,           // аргументы
    )
    if err != nil {
        log.Fatalf("Ошибка создания очереди: %s", err)
    }

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
    if err != nil {
        log.Fatalf("Ошибка подписки на очередь: %s", err)
    }

    // Обработка сообщений
    go func() {
        for d := range msgs {
            var order Order
            if err := json.Unmarshal(d.Body, &order); err != nil {
                log.Printf("Ошибка декодирования сообщения: %s", err)
                continue
            }
            log.Printf("Получено сообщение: %s", d.Body)
            sendEmail(order) // Отправка почты
        }
    }()

    log.Println("Ожидание сообщений...")
    select {}
}