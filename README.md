# 🔐 Auth System — авторизация по SMS и Email на Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Gin](https://img.shields.io/badge/Gin-1.9.1-00ADD8?style=flat&logo=gin)](https://gin-gonic.com/)

Готовая система аутентификации с поддержкой двухфакторной верификации через SMS и email.

## ✨ Возможности

- ✅ Регистрация и вход по email/телефону
- ✅ Подтверждение через SMS и email код
- ✅ JWT токены (access + refresh)
- ✅ Защита паролей bcrypt
- ✅ PostgreSQL + GORM
- ✅ Готов к Docker

## 🚀 Быстрый старт

```bash
# Клонирование
git clone https://github.com/yourusername/auth-system.git
cd auth-system

# Установка зависимостей
go mod download

# Настройка БД
createdb authdb

# Запуск
cp .env.example .env
# Отредактируйте .env

go run main.go
