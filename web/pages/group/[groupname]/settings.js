import { useState, useEffect } from 'react'
import { useRouter } from 'next/router'
import useSWR from 'swr'
import marked from 'marked'
import {sanitize} from 'dompurify'
const fetcher = (url) => fetch(url).then((res) => res.json())


export default function Index() {
  const router = useRouter()
  const {groupname} = router.query
  const [req, setReq] = useState({
    group_name: '',
    group_members: '',
  });
  useEffect(() => {
    if (router.asPath !== router.route) {
      fetcher(`/api/groups/${groupname}`).then((data) => {
        setReq({
          group_name: data.group_name,
          group_members: data.group_members,
        });
      })
    }
  }, [router])
  function handleSubmit(e) {
    e.preventDefault();
    //call api
    fetch(`/api/groups/${groupname}`, {
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
        <input type="submit" value="メンバー更新" />
      </form>
    </div>
  )
}
