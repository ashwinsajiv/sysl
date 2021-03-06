+++
# AUTOGENERATED BY byexample/generate.go
title= "Modules"
draft= false
description= ""
layout= "byexample"
weight = 8
topic = "Basics"
PlaygroundURL = "http://anz-bank.github.io/sysl-playground/?input=aW1wb3J0IGRlcHMKaW1wb3J0IG1vcmVkZXBzL21vcmVkZXBzCmltcG9ydCAvdGVzdHMvYmFuYW5hdHJlZQppbXBvcnQgLy9naXRodWIuY29tL2Fuei1iYW5rL3N5c2wtZXhhbXBsZXMvZGVtb3Mvc2ltcGxlL3NpbXBsZQppbXBvcnQgLy9naXRodWIuY29tL2Fuei1iYW5rL3N5c2wtZXhhbXBsZXMvZGVtb3Mvc2ltcGxlL3NpbXBsZUAwLjAuMQppbXBvcnQgZm9yZWlnbl9pbXBvcnRfc3dhZ2dlci55YW1sIGFzIGNvbS5mb28uYmFyLmFwcCB+c3dhZ2dlcgpNb2RlbFdpdGhEZXBzIFtwYWNrYWdlPSJtb2RlbCJdOgogICF0eXBlIFJlc3BvbnNlOgogICAgdXNlcklkIDw6IGludAogICAgaWQgPDogaW50CiAgICB0aXRsZSA8OiBzdHJpbmcKICAvcmVzcG9uc2VzOgogICAgL3tpZDw6aW50fToKICAgICAgR0VUOgogICAgICAgIERlcCA8LSBHRVQgL2RlcC97aWR9CiAgICAgICAgRGVwIDwtIEdFVCAvbW9yZWRlcC97aWR9CiAgICAgICAgRGVwMiA8LSBHRVQgL2RlcDIve2lkfQogICAgICAgIFRvZG9zIDwtIEdFVCAvdG9kb3Mve2lkfQogICAgICAgIEJhbmFuYXRyZWUgPC0gR0VUIC9iYW5hbmEve2lkfQogICAgICAgIHJldHVybiBSZXNwb25zZQo=&cmd="
GitRepoURL = "https://github.com/anz-bank/sysl/tree/master/demo/examples/Modules"
ID = "modules"
CodeWithoutComments = """import deps
import moredeps/moredeps
import /tests/bananatree
import //github.com/anz-bank/sysl-examples/demos/simple/simple
import //github.com/anz-bank/sysl-examples/demos/simple/simple@0.0.1
import foreign_import_swagger.yaml as com.foo.bar.app ~swagger
ModelWithDeps [package="model"]:
  !type Response:
    userId <: int
    id <: int
    title <: string
  /responses:
    /{id<:int}:
      GET:
        Dep <- GET /dep/{id}
        Dep <- GET /moredep/{id}
        Dep2 <- GET /dep2/{id}
        Todos <- GET /todos/{id}
        Bananatree <- GET /banana/{id}
        return Response
"""

Segs = [[
  
      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<pre class="chroma"><span class="nx">Dep</span> <span class="p">[</span><span class="kn">package</span><span class="p">=</span><span class="s">&#34;dep&#34;</span><span class="p">]:</span>
  <span class="p">!</span><span class="kd">type</span> <span class="nx">Dep</span><span class="p">:</span>
    <span class="nx">id</span> <span class="p">&lt;:</span> <span class="kt">int</span>
    <span class="nx">title</span> <span class="p">&lt;:</span> <span class="kt">string</span></pre>""",DocsRendered= """""",Image = ""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= true,CodeRendered="""<pre class="chroma">  <span class="o">/</span><span class="nx">dep</span><span class="p">:</span>
    <span class="o">/</span><span class="p">{</span><span class="nx">id</span><span class="p">&lt;:</span><span class="kt">int</span><span class="p">}:</span>
      <span class="nx">GET</span><span class="p">:</span>
        <span class="k">return</span> <span class="nx">Dep</span></pre>""",DocsRendered= """""",Image = ""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<pre class="chroma">  <span class="o">/</span><span class="nx">moredep</span><span class="p">:</span>
    <span class="o">/</span><span class="p">{</span><span class="nx">id</span><span class="p">&lt;:</span><span class="kt">int</span><span class="p">}:</span>
      <span class="nx">GET</span><span class="p">:</span>
        <span class="k">return</span> <span class="nx">Dep</span></pre>""",DocsRendered= """""",Image = ""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<pre class="chroma"><span class="nx">Dep2</span> <span class="p">[</span><span class="kn">package</span><span class="p">=</span><span class="s">&#34;dep2&#34;</span><span class="p">]:</span>
  <span class="p">!</span><span class="kd">type</span> <span class="nx">Dep2</span><span class="p">:</span>
    <span class="nx">id2</span> <span class="p">&lt;:</span> <span class="kt">int</span>
    <span class="nx">title2</span> <span class="p">&lt;:</span> <span class="kt">string</span></pre>""",DocsRendered= """""",Image = ""},

      {CodeEmpty= false,CodeLeading= false,CodeRun= false,CodeRendered="""<pre class="chroma">  <span class="o">/</span><span class="nx">dep2</span><span class="p">:</span>
    <span class="o">/</span><span class="p">{</span><span class="nx">id</span><span class="p">&lt;:</span><span class="kt">int</span><span class="p">}:</span>
      <span class="nx">GET</span><span class="p">:</span>
        <span class="k">return</span> <span class="nx">Dep2</span></pre>""",DocsRendered= """""",Image = ""},


],
[
  
      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<pre class="chroma"><span class="nx">swagger</span><span class="p">:</span> <span class="s">&#34;2.0&#34;</span>
<span class="nx">info</span><span class="p">:</span>
  <span class="nx">title</span><span class="p">:</span> <span class="nx">Simple</span>
<span class="nx">paths</span><span class="p">:</span>
  <span class="o">/</span><span class="nx">test</span><span class="p">:</span>
    <span class="nx">get</span><span class="p">:</span>
      <span class="nx">responses</span><span class="p">:</span>
        <span class="mi">200</span><span class="p">:</span>
          <span class="nx">description</span><span class="p">:</span> <span class="mi">200</span> <span class="nx">OK</span>
          <span class="nx">schema</span><span class="p">:</span>
            <span class="err">$</span><span class="nx">ref</span><span class="p">:</span> <span class="err">&#39;#</span><span class="o">/</span><span class="nx">definitions</span><span class="o">/</span><span class="nx">SimpleObj</span><span class="err">&#39;</span>
<span class="nx">definitions</span><span class="p">:</span>
  <span class="nx">SimpleObj</span><span class="p">:</span></pre>""",DocsRendered= """""",Image = ""},

      {CodeEmpty= false,CodeLeading= false,CodeRun= true,CodeRendered="""<pre class="chroma">    <span class="kd">type</span><span class="p">:</span> <span class="nx">object</span>
    <span class="nx">properties</span><span class="p">:</span>
      <span class="nx">name</span><span class="p">:</span>
        <span class="kd">type</span><span class="p">:</span> <span class="kt">string</span></pre>""",DocsRendered= """""",Image = ""},


],
[
  
      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<pre class="chroma">
<span class="kn">import</span> <span class="nx">deps</span></pre>""",DocsRendered= """<p>To import a sysl file from the same folder</p>
""",Image = ""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= true,CodeRendered="""<pre class="chroma">
<span class="kn">import</span> <span class="nx">moredeps</span><span class="o">/</span><span class="nx">moredeps</span></pre>""",DocsRendered= """<p>To import a sysl file from the subdirectory</p>
""",Image = ""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<pre class="chroma">
<span class="kn">import</span> <span class="o">/</span><span class="nx">tests</span><span class="o">/</span><span class="nx">bananatree</span></pre>""",DocsRendered= """<p>To import a sysl file via absolute path based on root</p>
""",Image = ""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<pre class="chroma">
<span class="kn">import</span> <span class="o">//</span><span class="nx">github</span><span class="p">.</span><span class="nx">com</span><span class="o">/</span><span class="nx">anz</span><span class="o">-</span><span class="nx">bank</span><span class="o">/</span><span class="nx">sysl</span><span class="o">-</span><span class="nx">examples</span><span class="o">/</span><span class="nx">demos</span><span class="o">/</span><span class="nx">simple</span><span class="o">/</span><span class="nx">simple</span></pre>""",DocsRendered= """<p>To import a sysl file from external repo. Add // at the beginning of the path</p>
""",Image = ""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<pre class="chroma">
<span class="kn">import</span> <span class="o">//</span><span class="nx">github</span><span class="p">.</span><span class="nx">com</span><span class="o">/</span><span class="nx">anz</span><span class="o">-</span><span class="nx">bank</span><span class="o">/</span><span class="nx">sysl</span><span class="o">-</span><span class="nx">examples</span><span class="o">/</span><span class="nx">demos</span><span class="o">/</span><span class="nx">simple</span><span class="o">/</span><span class="nx">simple</span><span class="err">@</span><span class="mf">0.0.1</span></pre>""",DocsRendered= """<p>To import a sysl file from external repo with fixed version</p>
""",Image = ""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<pre class="chroma">
<span class="kn">import</span> <span class="nx">foreign_import_swagger</span><span class="p">.</span><span class="nx">yaml</span> <span class="nx">as</span> <span class="nx">com</span><span class="p">.</span><span class="nx">foo</span><span class="p">.</span><span class="nx">bar</span><span class="p">.</span><span class="nx">app</span> <span class="err">~</span><span class="nx">swagger</span></pre>""",DocsRendered= """<p>To import a swagger file from the same folder. com.foo.bar.app is the alias of the module</p>
""",Image = ""},

      {CodeEmpty= false,CodeLeading= true,CodeRun= false,CodeRendered="""<pre class="chroma"><span class="nx">ModelWithDeps</span> <span class="p">[</span><span class="kn">package</span><span class="p">=</span><span class="s">&#34;model&#34;</span><span class="p">]:</span>
  <span class="p">!</span><span class="kd">type</span> <span class="nx">Response</span><span class="p">:</span>
    <span class="nx">userId</span> <span class="p">&lt;:</span> <span class="kt">int</span>
    <span class="nx">id</span> <span class="p">&lt;:</span> <span class="kt">int</span>
    <span class="nx">title</span> <span class="p">&lt;:</span> <span class="kt">string</span></pre>""",DocsRendered= """""",Image = ""},

      {CodeEmpty= false,CodeLeading= false,CodeRun= false,CodeRendered="""<pre class="chroma">  <span class="o">/</span><span class="nx">responses</span><span class="p">:</span>
    <span class="o">/</span><span class="p">{</span><span class="nx">id</span><span class="p">&lt;:</span><span class="kt">int</span><span class="p">}:</span>
      <span class="nx">GET</span><span class="p">:</span>
        <span class="nx">Dep</span> <span class="o">&lt;-</span> <span class="nx">GET</span> <span class="o">/</span><span class="nx">dep</span><span class="o">/</span><span class="p">{</span><span class="nx">id</span><span class="p">}</span>
        <span class="nx">Dep</span> <span class="o">&lt;-</span> <span class="nx">GET</span> <span class="o">/</span><span class="nx">moredep</span><span class="o">/</span><span class="p">{</span><span class="nx">id</span><span class="p">}</span>
        <span class="nx">Dep2</span> <span class="o">&lt;-</span> <span class="nx">GET</span> <span class="o">/</span><span class="nx">dep2</span><span class="o">/</span><span class="p">{</span><span class="nx">id</span><span class="p">}</span>
        <span class="nx">Todos</span> <span class="o">&lt;-</span> <span class="nx">GET</span> <span class="o">/</span><span class="nx">todos</span><span class="o">/</span><span class="p">{</span><span class="nx">id</span><span class="p">}</span>
        <span class="nx">Bananatree</span> <span class="o">&lt;-</span> <span class="nx">GET</span> <span class="o">/</span><span class="nx">banana</span><span class="o">/</span><span class="p">{</span><span class="nx">id</span><span class="p">}</span>
        <span class="k">return</span> <span class="nx">Response</span></pre>""",DocsRendered= """""",Image = ""},


],

]
+++


