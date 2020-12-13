import { useState, useEffect } from 'react'
import { useRouter } from 'next/router'
import useSWR from 'swr'
import marked from 'marked'
import {sanitize} from 'dompurify'
const fetcher = (url) => fetch(url).then((res) => res.json())


export default function Index() {
  const router = useRouter()
  const {article_id} = router.query
  const [mounted, setMounted] = useState(false)
  const [body, setBody] = useState('')
  useEffect(() => {
    setMounted(true)
    if (router.asPath !== router.route) {
      fetcher(`/api/articles/${article_id}`).then((data) => {
        setBody(data.body);
      })
    }
  }, [router])
  console.log(article_id)
  const { data, error } = useSWR(() => (mounted ? `/api/articles/${article_id}` : null), fetcher)
  if (error) return <div>Failed to load, Please Login</div>
  if (!data) return <div>Loading now...</div>
  let tags= []
  for(let tag of data.tags) {
    tags.push(<span key={`tag-${tag}`}>{tag} </span>)
  }
  function handleSubmit(e) {
    e.preventDefault();
    //call api
    fetch(`/api/articles/${article_id}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        'title': data.title,
        'article_path': data.article_path,
        'tags': data.tags,
        'group_name': data.group_name,
        'body': body,
      }),
    })
      .then((r) => {
        return r.json();
      })
      .then((data) => {
        router.push(`/article/${article_id}`);
      });
  }
  return (
    <div>
      <form onSubmit={handleSubmit}>
        <h1>タイトル: {`${data.title}`}</h1>
        <p>Path: {`${data.article_path}`}</p>
        <p>Tag: {tags}</p>
        <p>Group: {data.group_name}</p>
        <p>Article: </p>
        <textarea value={body} onChange={(e) => setBody(e.target.value)}/>
        <input type="submit" value="更新" />
      </form>
      <div dangerouslySetInnerHTML={{__html: sanitize(marked(body))}}/>
    </div>
  )
}
