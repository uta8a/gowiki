import { useState, useEffect } from 'react'
import { useRouter } from 'next/router'
import useSWR from 'swr'
import marked from 'marked'
import {sanitize} from 'dompurify'
const fetcher = (url) => fetch(url).then((res) => res.json())

const useMounted = () => {
  const [mounted, setMounted] = useState(false)
  useEffect(() => setMounted(true), [])
  return mounted
}

export default function Index() {
  const mounted = useMounted()
  const router = useRouter()
  const {id} = router.query
  const { data, error } = useSWR(() => (mounted ? `/api/articles/${id}` : null), fetcher)

  if (error) return <div>Failed to load, Please Login</div>
  if (!data) return <div>Loading now...</div>
  let body = sanitize(marked(data.body))
  let tags= []
  for(let tag of data.tags) {
  tags.push(<span key={`tag-${tag}`}>{tag} </span>)
  }
  return (
    <div>
      <h1>タイトル: {`${data.title}`}</h1>
      <p>Path: {`${data.article_path}`}</p>
      <p>Tag: {tags}</p>
      <p>Group: {data.group_name}</p>
      <p>Article: </p>
      <a href={`/edit/${id}`}>編集する</a>
      <div dangerouslySetInnerHTML={{__html: body}} />
    </div>
  )
}
