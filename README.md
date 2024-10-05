# MoodleAttendo

<p align="center">
<img src="https://github.com/user-attachments/assets/d527ab1e-f907-45b3-9712-135f07faaaee" alt="prbcare" width="400">
</p>

MoodleAttendo adalah aplikasi automation untuk menandai kehadiran pada website elearning berbasis Moodle. Aplikasi ini dibuat berdasarkan hasil dari reverse enginering pada salah satu website elearning berbasis moodle dengan versi 4.3.

## Environment Variables

MoodleAttendo membutuhkan environment variables berikut yang harus diatur sebelum dijalankan:

| **Kunci**    | **Tipe**     | **Deskripsi**                         | **Contoh**                                       |
|--------------|--------------|---------------------------------------|--------------------------------------------------|
| **HOSTNAME** | `string`     | Host website yang menggunakan Moodle. | `www.example.com`                                |
| **USERNAME** | `string`     | Username akun Moodle.                 | `username`                                       |
| **PASSWORD** | `string`     | Password akun Moodle.                 | `Admin#123`                                      |
| **TGCHAT**   | `string`     | ID chat Telegram.                     | `1214408099`                                     |
| **TGBOT**    | `string`     | Token bot Telegram.                   | `7286672841:AAGeF0rYJMixCJHEZ8P_7-_peaPYwJKw1rk` |

## Usage

### Docker Image

1. Fork repositori ini dan buat image menggunakan github actions atau lakukan pull image dari yang sudah ada di [sini](https://ghcr.io/scrkiddie/moodleattendo:latest).
2. Jalankan image dengan environment variables yang diperlukan dan argumen ID course yang bisa didapatkan di halaman course pada website berbasis moodle.

<p align="center">
<img src="https://github.com/user-attachments/assets/9c91fd0e-10a6-4a11-aa5f-59678741a3a2" alt="prbcare" width="400">
</p>

```shell
docker run -e HOSTNAME=$HOSTNAME -e USERNAME=$USERNAME -e PASSWORD=$PASSWORD -e TGCHAT=$TGCHAT -e TGBOT=$TGBOT image 1212
```

Untuk contoh penggunaan yang lebih lengkap lihat file `.circleci/config.yml`.

### Executable
1. Pastikan aplikasi chromium sudah terinstall. 
2. Clone repositori ini.
3. Build file executable.
4. Atur environment variables yang diperlukan.
5. Jalankan program dengan ID course sebagai argumen yang bisa didapatkan di halaman course pada website berbasis moodle.

<p align="center">
<img src="https://github.com/user-attachments/assets/9c91fd0e-10a6-4a11-aa5f-59678741a3a2" alt="prbcare" width="400">
</p>

Linux/macOS:
```shell
go build -o moodle_attendo cmd/moodle_attendo/main.go
export HOSTNAME=www.example.com
# Atur environment variables lainnya...
./moodle_attendo 1212
```

Windows:
```shell
go build -o moodle_attendo.exe cmd\moodle_attendo\main.go
set HOSTNAME=www.example.com
# Atur environment variables lainnya...
moodle_attendo.exe 1212
```


## Contribution
Jika Anda ingin berkontribusi pada pengembangan MoodleAttendo, silakan buat issue atau ajukan pull request di repositori ini.