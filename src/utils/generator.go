package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	mathRand "math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/twinj/uuid"
)

var (
	table              = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	alphaNumeric       = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	dateTimeZoneFormat = "2006-01-02 15:04 MST"
	dateFormat         = "2006-01-02"
	fileDateTimeFormat = "2006-01-02-15:04"
)

func GenerateUuid() string {
	newUuid := uuid.NewV4()

	return newUuid.String()
}

func GenerateRandomNumber(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func GenerateRandomAlphaNumberic(max int) string {
	b := make([]rune, max)
	for i := range b {
		b[i] = alphaNumeric[mathRand.Intn(len(alphaNumeric))]
	}
	return string(b)
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func GenerateRandomRangeNumber(min, max int) int {
	return mathRand.Intn(max-min) + min
}

func GenerateExpireOTP() time.Time {
	getDuration, _ := strconv.Atoi(os.Getenv("OTP_DURATION"))

	return time.Now().Add(time.Duration(getDuration) * time.Second)
}

func GenerateFormattedDate(getDate, getTime string) time.Time {
	var genDate time.Time

	getZone, _ := time.Now().Zone()

	if getTime != "" {
		combineDateTime := fmt.Sprintf("%s %s %s", getDate, getTime, getZone)
		genDate, _ = time.Parse(dateTimeZoneFormat, combineDateTime)

		return genDate.Local()
	}

	genDate, _ = time.Parse(dateFormat, getDate)

	return genDate.Local()
}

func GenerateDecimal(data int64) decimal.Decimal {
	genDecimal := decimal.NewFromInt(data)
	return genDecimal
}

func GenerateFormattedFilePath(file *multipart.FileHeader, userName string) string {
	currTimeUnique := time.Now().Local().Format(fileDateTimeFormat)

	newFileName := strings.Join([]string{userName, filepath.Ext(file.Filename)}, ".")

	filePath := fmt.Sprintf("/%s-%s", currTimeUnique, newFileName)

	return filePath
}

func GenerateFormattedFilePathHash(file *multipart.FileHeader) string {
	currTimeUnique := time.Now().Local().Format(fileDateTimeFormat)
	userName := GenerateSecureToken(50)

	newFileName := strings.Join([]string{userName, filepath.Ext(file.Filename)}, "")

	filePath := fmt.Sprintf("%s-%s", currTimeUnique, newFileName)

	return filePath
}

func GenerateValidPhoneNumber(phoneNumber string) string {
	regex := regexp.MustCompile(`^0+`)
	phoneNumber = regex.ReplaceAllString(phoneNumber, `+62`)

	regex = regexp.MustCompile(`^8+`)
	isMatch := regex.MatchString(phoneNumber)
	if isMatch {
		phoneNumber = "+62" + phoneNumber
	}

	return phoneNumber
}

func GenerateCodeVoucher(username, existedCodeVoucher string) string {
	getRandAlphaNum := GenerateRandomAlphaNumberic(4)
	splitUsername := strings.Split(username, "")
	getFourDigitUsername := make([]string, len(splitUsername))

	for i := 0; i < 4; i++ {
		getFourDigitUsername = append(getFourDigitUsername, splitUsername[i])
	}

	result := strings.ToUpper(strings.Join(getFourDigitUsername, "")) + "_" + strings.ToUpper(getRandAlphaNum)

	if result != existedCodeVoucher || existedCodeVoucher == "" {
		return result
	}

	return GenerateCodeVoucher(username, existedCodeVoucher)
}

func GenerateRomanMonth(num int) string {
	values := []int{
		1000, 900, 500, 400,
		100, 90, 50, 40,
		10, 9, 5, 4, 1,
	}

	symbols := []string{
		"M", "CM", "D", "CD",
		"C", "XC", "L", "XL",
		"X", "IX", "V", "IV",
		"I"}
	roman := ""
	i := 0

	for num > 0 {
		// calculate the number of times this num is completely divisible by values[i]
		// times will only be > 0, when num >= values[i]
		k := num / values[i]
		for j := 0; j < k; j++ {
			//buildup roman numeral
			roman += symbols[i]

			//reduce the value of num.
			num -= values[i]
		}
		i++
	}
	return roman
}

func GenerateSchemeAndHost(ctx *gin.Context) (host string, scheme string) {
	host = ctx.Request.Host
	scheme = "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	return host, scheme
}

func GenerateNoOrder(lastInvoiceData string) string {
	getSameNumFromInvoice := strings.Split(lastInvoiceData, "/")
	getCurrTime := time.Now().Local()
	getTwoDigitYear := getCurrTime.Format("06")

	_, month, _ := getCurrTime.Date()
	getRomanMonth := GenerateRomanMonth(int(month))

	return fmt.Sprintf("%s/KLIPZ/ORD/%s/%s", getSameNumFromInvoice[0], getRomanMonth, getTwoDigitYear)
}

func GenerateInvoiceNumber() string {
	getCurrTime := time.Now().Local()
	getTwoDigitYear := getCurrTime.Format("06")

	_, month, _ := getCurrTime.Date()
	getRomanMonth := GenerateRomanMonth(int(month))

	return fmt.Sprintf("INV/%s/%s/%s", getRomanMonth, getTwoDigitYear, GenerateRandomNumber(4))
}
