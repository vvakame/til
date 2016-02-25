# NOTE

実行ログ
cat ~/Library/Logs/Unity/Player.log

設定ファイル
cat ~/Library/Preferences/unity.vvakame.webcamtexture.plist

設定ファイル内容確認
plutil -p ~/Library/Preferences/unity.vvakame.webcamtexture.plist

設定内容更新
plutil -replace CameraName -string "FaceTime HD Camera" ~/Library/Preferences/unity.vvakame.webcamtexture.plist

で、設定変更できてほしいけど全然できない。つらい。
