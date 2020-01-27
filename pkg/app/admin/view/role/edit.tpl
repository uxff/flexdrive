{{ define "role/edit.tpl" }}

{{template "common/head.tpl" .}}
{{template "common/partheader.tpl" .}}

<div class="container">
    <div class="row">
        <ul class="breadcrumb">
          <li><a href="/">首页</a></li>
          <li><a href="/role/list">角色</a></li>
          <li class="active">编辑角色</li>
        </ul>
    	<div class="col-md-8 col-md-offset-2">
    		<div class="panel panel-default">
			  	<div class="panel-heading text-center">
			    	<h3 class="panel-title"><strong>编辑角色</strong></h3>
			 	</div> 

            <div class="panel-body">
                <form accept-charset="utf-8" role="form" class="form-horizontal" method="POST" action="/role/edit/{{.RoleEnt.Id}}">
                    <div class="form-group">
                        <label for="" class="col-sm-3 control-label">角色id</label>
                        <div class="col-sm-8 control-label" style="text-align: left">
                            {{.RoleEnt.Id}}
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="inputName" class="col-sm-3 control-label">名称</label>
                        <div class="col-sm-8">
                            <input class="form-control" placeholder="例: 运营" name="name" value="{{.RoleEnt.Name}}" type="input" required 
                                    id="inputName" />
                        </div>
                    </div>
                    <!--下面编辑权限-->
                    <div class="col-sm-12">
                        {{range $accessGroup := .allAccessItems}}

                        <div class="panel panel-default">
                            <div class="panel-heading ">
                                <span class="panel-title">{{$accessGroup.Name}}</span>
                            </div>
                            <div class="panel-body">
                                {{range $accessItem := $accessGroup.Sub}}
                                    <label for="checkbox-{{$accessItem.PermitRoute}}" class="col-sm-4 ">
                                        <input id="checkbox-{{$accessItem.PermitRoute}}" name="accessRoute[{{$accessItem.PermitRoute}}]" value="1" type="checkbox" {{if $accessItem.Access}}checked{{end}} />{{$accessItem.Name}}
                                    </label>
                                {{end}}
                            </div>
                        </div>
                        {{end}}
                    </div>
                    
                    <div class="form-group text-center">
                        <div class="col-sm-5"></div>
                        <div class="col-sm-2">
                            <input class="btn btn-success btn-block" type="submit" value="提交">
                        </div>
                        <div class="col-sm-5"></div>
                    </div>
            </form>
            </div>
			</div>
		</div>
	</div>
</div>

{{template "common/partfooter.tpl"}}
{{template "common/foot.tpl"}}

{{ end }}