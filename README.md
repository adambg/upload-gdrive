# upload-gdrive
HTTP listener to receive file by POST and submit directly to Google Drive

# Configuration
* Create credentials.json: Follow the quickstart steps: https://developers.google.com/gmail/api/quickstart/go
* Set `googleDriveID` - you can get it from the URL of the drive folder location where you wish the files will be uploaded

# How to use
using cURL:
```
curl -i -X POST -H "Content-Type: multipart/form-data" -F "file=@myfile.doc;filename=myfile.doc" http://server/upload/
```

using Java:
```
try {
    FileInputStream fstrm = new FileInputStream("myfile.doc");
    HttpFileUpload hfu = new HttpFileUpload("http://server/upload/" + "myfile.doc", "myfile.doc", "myfile.doc");
    hfu.Send_Now(fstrm);
} catch (Exception e) {
    // error
}
```
