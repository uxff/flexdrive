{{if gt .paginator.PageNums 1}}
<ul class="pagination">
{{if .paginator.HasPrev}}
    <li><a href="{{.paginator.PageLinkFirst}}" title="首页"><i class="glyphicon glyphicon-step-backward"></i></a></li>
    <li><a href="{{.paginator.PageLinkPrev}}"><i class="glyphicon glyphicon-backward"></i></a></li>
{{else}}
    <li class="disabled"><a title="首页"><i class="glyphicon glyphicon-step-backward"></i></a></li>
    <li class="disabled"><a><i class="glyphicon glyphicon-backward"></i></a></li>
{{end}}
{{range $index, $page := .paginator.Pages}}
    <li{{if $.paginator.IsActive .}} class="active"{{end}}>
        <a href="{{$.paginator.PageLink $page}}">{{$page}}</a>
    </li>
{{end}}
{{if .paginator.HasNext}}
    <li><a href="{{.paginator.PageLinkNext}}"><i class="glyphicon glyphicon-forward"></i></a></li>
    <li><a href="{{.paginator.PageLinkLast}}" title="末页"><i class="glyphicon glyphicon-step-forward"></i></a></li>
{{else}}
    <li class="disabled"><a><i class="glyphicon glyphicon-forward"></i></a></li>
    <li class="disabled"><a title="末页"><i class="glyphicon glyphicon-step-forward"></i></a></li>
{{end}}
    <li class="disabled"><a>共{{.paginator.PageNums}}页</a></li>
</ul>
{{end}}