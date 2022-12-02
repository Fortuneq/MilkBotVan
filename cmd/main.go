package main

import (
	"github.com/jmoiron/sqlx"
	tele "gopkg.in/telebot.v3"
	"log"
	"milk/config"
	"milk/pkg/logger"
	"milk/pkg/postgres"
	"milk/pkg/telebotCalendar"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfgFile, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewAPILogger(cfg)
	err = appLogger.InitLogger()
	if err != nil {
		log.Fatalf("Cannot init logger: %v", err.Error())
	}
	appLogger.Infof("logger successfully started with - Level: %s", cfg.Logger.Level)

	psqlDB, err := postgres.InitPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("PostgreSQL init error: %s", err)
	} else {
		appLogger.Infof("PostgreSQL connected, status: %#v", psqlDB.Stats())
	}
	defer func(psqlDB *sqlx.DB) {
		err = psqlDB.Close()
		if err != nil {
			appLogger.Infof(err.Error())
		} else {
			appLogger.Info("PostgreSQL closed properly")
		}
	}(psqlDB)

	prefMilk := tele.Settings{
		Token:  cfg.Telegram.VerifySystem.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	MilkBot, err := tele.NewBot(prefMilk)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		MilkBot.Start()
	}()

	selector := &tele.ReplyMarkup{}

	btnBook := selector.Text("Записаться на прием")
	btnCancel := selector.Text("Отменить запись")
	//selector.Reply(
	//	selector.Row(btnBook),
	//	selector.Row(btnCancel),
	//)
	selector.InlineKeyboard = telebotCalendar.GenerateCalendar(2022, 1)
	MilkBot.Handle(&btnBook, func(c tele.Context) error {
		return c.EditOrReply("Выберите доступные дату и время", selector)
	})

	MilkBot.Handle(&btnCancel, func(c tele.Context) error {
		return c.EditOrReply("Список ваших записей")
	})

	MilkBot.Handle("/start", func(c tele.Context) error {
		c.Send("Добрый день, через этого бота вы можете записаться на услугу в Имидж стидии VAN", selector)
		if err != nil {
			return err
		}
		return nil
	})
	appLogger.Info("Application has started")

	exit := make(chan os.Signal, 1)

	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-exit

	appLogger.Info("Application has been shut down")
}
