import { useState, useEffect } from 'react'
import { useRouter } from 'next/router'
import useSWR from 'swr'
import marked from 'marked'
import {sanitize} from 'dompurify'
const fetcher = (url) => fetch(url).then((res) => res.json())


export default function Index() {
  const router = useRouter()
  const [req, setReq] = useState({
    group_name: '',
    group_members: '',
  });
  function handleSubmit(e) {
    e.preventDefault();
    //call api
    fetch(`/api/groups`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        'group_name': req.group_name,
        'group_members': req.group_members.split(',').map(v => v.trim()),
      }),
    })
      .then((r) => {
        return r.json();
      })
      .then((data) => {
        router.push(`/group/${data.group_name}/settings`);
      });
  };
  return (
    <div>
      <form onSubmit={handleSubmit}>
        <p>Group: <input type="text" value={req.group_name} onChange={(e) => setReq({...req, ...{group_name: e.target.value}})}/></p>
        <p>GroupMembers: </p>
        <textarea value={req.group_members} onChange={(e) => setReq({...req, ...{group_members: e.target.value}})}/>
        <input type="submit" value="新規作成" />
      </form>
    </div>
  )
}
