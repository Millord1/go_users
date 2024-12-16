package utils

import (
	"bufio"
	"fmt"
	"log"
	"microservices/models"
	"os"
	"time"

	"github.com/mdp/qrterminal/v3"
	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
)

func GenerateTOTPWithSecret(user *models.User, randSecret string) {
	totp := gotp.NewDefaultTOTP(randSecret)

	user.Totp = randSecret
	_, err := user.HashTotp()
	if err != nil {
		log.Fatalln(err)
	}

	uri := totp.ProvisioningUri(user.Email, "Millord")
	fmt.Println("secret key: ", uri)

	qrcode.WriteFile(uri, qrcode.Medium, 256, "./assets/"+user.Username+".png")

	qrterminal.GenerateWithConfig(uri, qrterminal.Config{
		Level:     qrterminal.L,
		Writer:    os.Stdout,
		BlackChar: qrterminal.BLACK,
		WhiteChar: qrterminal.WHITE,
	})

	fmt.Println("\n scan it")
}

func VerifyOtp(randSecret string, otp string) bool {
	totp := gotp.NewDefaultTOTP(randSecret)
	return totp.Verify(otp, time.Now().Unix())
}

func DevVerifyOTP(randSecret string) {
	// Usefull to debug from terminal
	totp := gotp.NewDefaultTOTP(randSecret)

	fmt.Printf("Please enter your OTP from your app: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userInput := scanner.Text()

	if totp.Verify(userInput, time.Now().Unix()) {
		fmt.Println("Successful !!!")
	} else {
		fmt.Println("Failed :(")
	}
}
