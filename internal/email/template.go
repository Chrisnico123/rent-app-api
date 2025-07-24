package email

import "fmt"

func OtpEmailTemplate(code string) string {
	return fmt.Sprintf(`
	<html>
		<body>
			<h2>Kode OTP Anda</h2>
			<p>Masukkan kode berikut untuk verifikasi: <b style="font-size:2em;">%s</b></p>
			<p>Kode ini hanya berlaku 4 digit dan akan kadaluarsa dalam 7 hari.</p>
		</body>
	</html>
	`, code)
}
