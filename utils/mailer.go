package utils

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(toEmail, token string) error {

	prod := os.Getenv("APP_ENV")
	var link string

	if prod != "production" {
		link = fmt.Sprintf("http://localhost:8888/change-password?token=%s", token)
	} else {
		link = fmt.Sprintf("https://gocatat.my.id/change-password?token=%s", token)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Verifikasi Akun GoCatat")

	// Template HTML Profesional & Responsive (Bagus di HP & Desktop)
	body := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Verifikasi Akun GoCatat</title>
		<style>
			body { font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; background-color: #f4f6f8; margin: 0; padding: 0; -webkit-font-smoothing: antialiased; }
			.container { max-width: 600px; margin: 40px auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.05); }
			.header { background-color: #2563eb; padding: 30px; text-align: center; }
			.header h1 { color: #ffffff; margin: 0; font-size: 24px; font-weight: bold; letter-spacing: 0.5px; }
			.content { padding: 40px 30px; color: #334155; line-height: 1.6; }
			.content h2 { color: #1e293b; font-size: 20px; margin-top: 0; }
			.content p { font-size: 16px; margin-bottom: 24px; }
			.btn-wrapper { text-align: center; margin: 35px 0; }
			.btn { background-color: #2563eb; color: #ffffff !important; text-decoration: none; padding: 14px 30px; font-size: 16px; font-weight: 600; border-radius: 6px; display: inline-block; transition: background-color 0.2s; box-shadow: 0 2px 5px rgba(37,99,235,0.2); }
			.footer { background-color: #f8fafc; padding: 20px 30px; text-align: center; font-size: 13px; color: #94a3b8; border-top: 1px solid #f1f5f9; }
			.footer a { color: #64748b; text-decoration: underline; }
			.fallback-link { font-size: 12px; color: #64748b; word-break: break-all; margin-top: 20px; padding: 15px; background-color: #f8fafc; border-radius: 4px; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>GoCatat</h1>
			</div>
			<div class="content">
				<h2>Halo,</h2>
				<p>Terima kasih telah mendaftar di <strong>GoCatat</strong>. Langkah terakhir untuk mengaktifkan akun Anda adalah dengan melakukan verifikasi alamat email.</p>
				
				<div class="btn-wrapper">
					<a href="%s" class="btn" target="_blank">Verifikasi Akun Saya</a>
				</div>
				
				<p>Tautan verifikasi ini hanya berlaku selama 24 jam. Jika Anda tidak merasa melakukan pendaftaran ini, Anda dapat mengabaikan email ini dengan aman.</p>
				
				<div class="fallback-link">
					<p style="margin: 0 0 5px 0; font-weight: bold;">Tombol tidak berfungsi? Salin link di bawah ke browser Anda:</p>
					<a href="%s">%s</a>
				</div>
			</div>
			<div class="footer">
				<p>&copy; 2026 GoCatat. Hak Cipta Dilindungi.</p>
				<p>Ini adalah email otomatis, mohon tidak membalas email ini.</p>
			</div>
		</div>
	</body>
	</html>
	`, link, link, link)

	m.SetBody("text/html", body)

	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"))

	return d.DialAndSend(m)
}

func SendChangePasswordEmail(toEmail, token string) error {

	prod := os.Getenv("APP_ENV")
	var link string

	if prod != "production" {
		link = fmt.Sprintf("http://localhost:8888/change-password?token=%s", token)
	} else {
		link = fmt.Sprintf("https://gocatat.my.id/change-password?token=%s", token)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Change Password Akun GoCatat")
    body := fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Ubah Password Akun GoCatat</title>
        <style>
            body { font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; background-color: #f4f6f8; margin: 0; padding: 0; -webkit-font-smoothing: antialiased; }
            .container { max-width: 600px; margin: 40px auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.05); }
            .header { background-color: #2563eb; padding: 30px; text-align: center; }
            .header h1 { color: #ffffff; margin: 0; font-size: 24px; font-weight: bold; letter-spacing: 0.5px; }
            .content { padding: 40px 30px; color: #334155; line-height: 1.6; }
            .content h2 { color: #1e293b; font-size: 20px; margin-top: 0; }
            .content p { font-size: 16px; margin-bottom: 24px; }
            .btn-wrapper { text-align: center; margin: 35px 0; }
            .btn { background-color: #2563eb; color: #ffffff !important; text-decoration: none; padding: 14px 30px; font-size: 16px; font-weight: 600; border-radius: 6px; display: inline-block; transition: background-color 0.2s; box-shadow: 0 2px 5px rgba(37,99,235,0.2); }
            .footer { background-color: #f8fafc; padding: 20px 30px; text-align: center; font-size: 13px; color: #94a3b8; border-top: 1px solid #f1f5f9; }
            .footer a { color: #64748b; text-decoration: underline; }
            .fallback-link { font-size: 12px; color: #64748b; word-break: break-all; margin-top: 20px; padding: 15px; background-color: #f8fafc; border-radius: 4px; }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>GoCatat</h1>
            </div>
            <div class="content">
                <h2>Halo,</h2>
                <p>Kami menerima permintaan untuk mengubah password akun <strong>GoCatat</strong> Anda. Untuk mengatur ulang password Anda, silakan klik tombol di bawah ini.</p>
                
                <div class="btn-wrapper">
                    <a href="%s" class="btn" target="_blank">Ubah Password Saya</a>
                </div>
                
                <p>Tautan ini hanya berlaku selama 24 jam. Jika Anda tidak merasa melakukan permintaan perubahan password ini, Anda dapat mengabaikan email ini dengan aman dan password Anda tidak akan diubah.</p>
                
                <div class="fallback-link">
                    <p style="margin: 0 0 5px 0; font-weight: bold;">Tombol tidak berfungsi? Salin link di bawah ke browser Anda:</p>
                    <a href="%s">%s</a>
                </div>
            </div>
            <div class="footer">
                <p>&copy; 2026 GoCatat. Hak Cipta Dilindungi.</p>
                <p>Ini adalah email otomatis, mohon tidak membalas email ini.</p>
            </div>
        </div>
    </body>
    </html>
    `, link, link, link)

    m.SetBody("text/html", body)

    port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
    d := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"))

    return d.DialAndSend(m)
}
