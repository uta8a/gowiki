import { useState, useEffect } from 'react'
import { useRouter } from 'next/router'
import useSWR from 'swr'

const fetcher = (url) => fetch(url).then((res) => res.json())

const useMounted = () => {
  const [mounted, setMounted] = useState(false)
  useEffect(() => setMounted(true), [])
  return mounted
}

export default function Index() {
  const mounted = useMounted()
  const router = useRouter()
  const { data, error } = useSWR(() => (mounted ? '/api/articles' : null), fetcher)

  if (error) return <div>Failed to load, Please Login</div>
  if (!data) return <div>Loading now...</div>
  console.log(data.groups)
  let list = []
  for(let g of (data.groups ? data.groups : [])) {
    let list_child = []
    for(let id of (g.articles_id ? g.articles_id : [])){
    list_child.push(<p key={`article-key-${id}`}>- <a href={`/article/${id}`}>{`article: ${id}`}</a></p>)
    }
    list.push(
      <div key={`key-div-${g.group_name}`}>
        <h2>Group: <a href={`/group/${g.group_name}/settings`}>{g.group_name}</a></h2>
        {list_child}
      </div>
    )
  }
  console.log("list",list)
  return (
    <div>
      <h1>User Page</h1>
      <div>
        <a href="/new">記事を新しく作る</a>
      </div>
      <div>
        <span>まずはグループを作ってみましょう！ → </span>
        <a href="/group">グループを新しく作る</a>
      </div>
      <div>
      {list}
      </div>
    </div>
  )
}
