import { useState, useEffect } from 'react'
import { useRouter } from 'next/router'
import useSWR from 'swr'
import marked from 'marked'
import {sanitize} from 'dompurify'
const fetcher = (url) => fetch(url).then((res) => res.json())


export default function Index() {
  const router = useRouter()
  const [req, setReq] = useState({
    title: '',
    article_path: '',
    tags: [''],
    group_name: '',
    body: ''
  });
  function handleSubmit(e) {
    e.preventDefault();
    //call api
    fetch(`/api/articles`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        'title': req.title,
        'article_path': req.article_path,
        'tags': req.tags,
        'group_name': req.group_name,
        'body': req.body,
      }),
    })
      .then((r) => {
        return r.json();
      })
      .then((data) => {
        router.push(`/article/${data.article_id}`);
      });
  };
  return (
    <div>
      <form onSubmit={handleSubmit}>
        <h1>タイトル: <input type="text" value={req.title} onChange={(e) => setReq({...req, ...{title: e.target.value}})}/></h1>
        <p>Path: <input type="text" value={req.article_path} onChange={(e) => setReq({...req, ...{article_path: e.target.value}})}/></p>
        <p>Group: <input type="text" value={req.group_name} onChange={(e) => setReq({...req, ...{group_name: e.target.value}})}/></p>
        <p>Article: </p>
        <textarea value={req.body} onChange={(e) => setReq({...req, ...{body: e.target.value}})}/>
        <input type="submit" value="新規作成" />
      </form>
      <div dangerouslySetInnerHTML={{__html: marked(req.body ? req.body: '')}}/>
    </div>
  )
}
