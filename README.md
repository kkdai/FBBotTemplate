FBBotTemplate: A simple Golang Facebook Bot Template 
==============

[![Join the chat at https://gitter.im/kkdai/FBBotTemplate](https://badges.gitter.im/kkdai/FBBotTemplate.svg)](https://gitter.im/kkdai/FBBotTemplate?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

 [![GoDoc](https://godoc.org/github.com/kkdai/FBBotTemplate.svg?status.svg)](https://godoc.org/github.com/kkdai/FBBotTemplate.svg)  [![Build Status](https://travis-ci.org/kkdai/FBBotTemplate.svg?branch=master)](https://travis-ci.org/kkdai/FBBotTemplate.svg)
 
 [![](https://goreportcard.com/badge/github.com/kkdai/FBBotTemplate)](https://goreportcard.com/badge/github.com/kkdai/FBBotTemplate)



Installation and Usage
=============

### 1. Create a Facebook Page

Make sure you already have Facebook account, if you need use FB Bot.

### 2. Just Deploy the same on Heroku

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

Remember your heroku, ID.

<br><br>

### 3. Enable Fixed IP service

![](images/linebotH2.png)
Clone the heroku git to your locally, use follow command to setup ([Fixie](https://elements.heroku.com/addons/fixie)) service for free.

#### There are two ways to add this "add-on" for heroku.

![](images/linebotH1.png)

1. Through heroku dashboard:
	- Launch Herou [dashboard](https://dashboard.heroku.com)
	- Go your deploy app page
	- Go to "Resource"
	- Go to "FIND MORE ADD-ONS" 
	- Find "Fixie"

2. Through [Heroku Toolbelt](https://toolbelt.heroku.com/)
	- `$ heroku login`
	- `$ heroku git:clone -a <YOUR_HEROKU_APP_ID>`
	- `$ cd linebotkkdaitest`
	- `$ heroku addons:create fixie:tricycle`

Remember your IP information. 


### 4. Back to Line Bot Dashboard, setup basic API

Setup your basic account information. Here is some info you will need to know.

- `Callback URL`: https://{YOUR_HEROKU_SERVER_ID}.herokuapp.com:443/callback

Go to `Server IP White List`, fill the IP from [Fixie](https://elements.heroku.com/addons/fixie)

You will get following info, need fill back to Heroku.

- Channel ID
- Channel Secret
- MID

### 5. Back to Heroku again to setup environment variables

- Go to dashboard
- Go to "Setting"
- Go to "Config Variables", add following variables:
	- "ChannelID"
	- "ChannelSecret"
	- "MID"

It all done.	



### Chinese Tutorial:

如果你看得懂繁體中文，這裡有[中文的介紹](http://www.evanlin.com/create-your-line-bot-golang/) 

Inspired By
=============

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

