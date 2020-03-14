{{ define "share/detail.tpl" }}
{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}
<div class="container">

    <div class="row">
        <ul class="breadcrumb">
            <li><a href="/">首页</a></li>
            <li class="active">分享详情</li>
        </ul>
    </div>

    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-8" >
            <div class="panel panel-success">
                <div class="panel-heading">
                    <h3 class="text-center">分享id： {{.ShareItem.Id}}</h3>
                </div>
                <div class="panel-body">
                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-4" >
                            分享时间
                        </div>
                        <div class="col-md-4" >
                            {{.ShareItem.Created}}
                        </div>
                        <div class="col-md-2"></div>
                    </div>
                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-4" >
                            分享人id
                        </div>
                        <div class="col-md-4" >
                            {{.ShareItem.UserId}}
                        </div>
                        <div class="col-md-2"></div>
                    </div>
                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-4" >
                            文件名
                        </div>
                        <div class="col-md-4" >
                            {{.ShareItem.FileName}}
                        </div>
                        <div class="col-md-2"></div>
                    </div>
                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-4" >
                            大小
                        </div>
                        <div class="col-md-4" >
                            {{size4Human .ShareItem.UserFile.Size}}
                        </div>
                        <div class="col-md-2"></div>
                    </div>
                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-4" >
                            有效期
                        </div>
                        <div class="col-md-4" >
                            {{.ShareItem.Expired}}
                        </div>
                        <div class="col-md-2"></div>
                    </div>
                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-4" >
                            下载地址
                        </div>
                        <div class="col-md-4" >
                            <a href="/file/{{.ShareItem.FileHash}}/{{.ShareItem.FileName}}" {{.ShareItem.Status}}>下载</a>
                        </div>
                        <div class="col-md-2"></div>
                    </div>

                </div>
            </div>
        </div>
        <div class="col-md-2"></div>
    </div>
</div>
{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}
{{end}}