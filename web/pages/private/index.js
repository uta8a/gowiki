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
  const { data, error } = useSWR(() => (mounted ? '/api/privatecheck' : null), fetcher)

  if (error) return <div>Failed to load, Please Login</div>
  if (!data) return <div>Loading now...</div>

  return (
    <div>
      <h1>Private Page</h1>
      <p>status code: {data.status}</p>
      <p>message: {data.message}</p>
    </div>
  )
}
