# 阿里云oss前端直传
> object store directly 

## 安装

  ```bash
  go get -u github.com/soonio/osd
  ```

## 获取上传配置并上传
  ```go
	var filename = "zhuye.png"
	
	var proxy = New(
		os.Getenv("key"),
		os.Getenv("secret"),
		os.Getenv("host"),
		UsePrefix("/local/"),
		UseDuration(120),
		UseCallback("https://www.iosoon.cn/api/any"),
	)
	var conf = proxy.Signature("/var/www/")

	url := conf.Host
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("key", conf.Dir+filename)
	_ = writer.WriteField("OSSAccessKeyId", conf.AccessId)
	_ = writer.WriteField("policy", conf.Policy)
	_ = writer.WriteField("Signature", conf.Signature)
	_ = writer.WriteField("callback", conf.Callback)
	_ = writer.WriteField("success_action_status", "200")

	file, err := os.Open("/local_dir/" + filename)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	part5, err := writer.CreateFormFile("file", filename)
	_, err = io.Copy(part5, file)
	if err != nil {
		panic(err)
	}

	err = writer.Close()
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
  ```

## 验证回调请求

  ```go
  r, err := http.NewRequest(
    "POST",
    "https://www.iosoon.cn/api/any",
    bytes.NewBufferString(`filename="local/var/www/zhuye.png"&size=5385&mimeType="application/octet-stream"&height=256&width=256`),
  )
  if err != nil {
    panic(err)
  }
  
  r.Header.Set("Authorization", "E6i//T0rnM/xu4wQjnKODWSYsyQdR2e8KIreuRihpBQzoy70q7KmdvrIR3vxNGi8OdCswE+QOfPt+O3PvA5NMA==")
  r.Header.Set("Content-Md5", "nbXlrGfwFyDIgSJjs1PBHA==")
  r.Header.Set("X-Oss-Pub-Key-Url", "aHR0cHM6Ly9nb3NzcHVibGljLmFsaWNkbi5jb20vY2FsbGJhY2tfcHViX2tleV92MS5wZW0=")
  
  var res = Verify(r)
  fmt.Println(res)
  ```