FBBotTemplate: A simple Golang Facebook Bot Template 
==============

[![Join the chat at https://gitter.im/kkdai/FBBotTemplate](https://badges.gitter.im/kkdai/FBBotTemplate.svg)](https://gitter.im/kkdai/FBBotTemplate?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

 [![GoDoc](https://godoc.org/github.com/kkdai/FBBotTemplate.svg?status.svg)](https://godoc.org/github.com/kkdai/FBBotTemplate.svg)  [![Build Status](https://travis-ci.org/kkdai/FBBotTemplate.svg?branch=master)](https://travis-ci.org/kkdai/FBBotTemplate.svg)
 
 [![](https://goreportcard.com/badge/github.com/kkdai/FBBotTemplate)](https://goreportcard.com/report/github.com/kkdai/FBBotTemplate)



Installation and Usage
=============


### 1. Just Deploy the same on Heroku

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

Remember your heroku ID and app address. ex: `https://APP_ADDRESS.herokuapp.com/`

### 2. Create a Facebook Page

Make sure you already have Facebook account, if you need use FB Bot.


![](images/Bot1.png)

### 3. Go to Facebook Developer to create App

Create App, need select as follow:

- New App type [Web App], create app
- Add new product [Messenger]

![](images/Bot2.png)


### 4. Configuration Messenger Bot

Get token from Faccebook page.

![](images/Bot4.png)

- Select a "Page" you own.
- Go to "Meseenger" product.
- It will generate `token`
- copy it and store it.

### 5. Paste Token to Heroku

![](images/Bot5.png)

Go to heroku dashboard, go to "Setting" -> "Config Variables".

- Add "Config Vars"
- Name -> "TOKEN"
- Value use  `token` facebook app.

### 6. Back Facebook App configuration

![](images/Bot6.png)

- Go to "Messenger" product
- Go to "Setup Webhooks"
- Fill `https://APP_ADDRESS.herokuapp.com/webhook` in callback URL.
- Fill your `token` in "token".
- Checked the checkbox "message_deliveres", "messages".

If your configuration is correct, it will show "completed".

![](images/Bot7.png)

Also remember to choose correct page in this setting.


### 7. Request basic permissions

![](images/Bot8.png)
 
In Messenger application review, press "Request Premission".
 
- Checked the "pages_messaging".



## How to testing it.

![](images/Bot9.png)

- Go to your spcific "page" in Facebook.
- Press "Send Message"

![](images/Bot10.png)


## How to modification your Bot Code

- Download code from Heroku `git clone https://git.heroku.com/APP_ADDRESS.git`
- Modify code on `main.go` first. especially in `MessageReceived()`.
- Commit and push it back to heroku `git push heroku master`.


Inspired By
=============

- [Provides a GO SDK for Facebook's messenger-platform](https://github.com/maciekmm/messenger-platform-go-sdk)
- [FB: Messenger Platform/Getting Started](https://developers.facebook.com/docs/messenger-platform/quickstart/)
- [[教學] Facebook Messenger API](http://huli.logdown.com/posts/709641-teaching-facebook-messenger-api)
- [用Python開發Facebook Bot](https://medium.com/dualcores-studio/%E7%94%A8python%E9%96%8B%E7%99%BCfacebook-bot-26594f13f9f7#.bunklnnue)



License
---------------

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

