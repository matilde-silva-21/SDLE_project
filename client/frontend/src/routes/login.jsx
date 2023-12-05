import { useEffect, useState } from "react"

export default function LoginPage() {
    const [username, setUsername] = useState("")
    
    useEffect(() => {
        const cookies = document.cookie.split(';')
        console.log(cookies)
        const session = cookies.find(cookie => cookie.split('=')[0].trim() === "session")
        console.log(session)
        
        if (session !== undefined) {
            const cookie = session.split('=')[1]

            if (cookie !== null) {
                document.location = '/'
            }
        }
    }, [])

    async function login(event) {
        event.preventDefault()

        const login = await fetch("http://localhost:8080/login", {
            method: 'POST',
            body: JSON.stringify({'username': `${username}`}),
            mode: 'cors',
            credentials: "include",
            headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        }})

        if (login.ok) {
            document.location = "/"
        }
    }

    return (
        <>
            <div className="container bg-pink-100 rounded-md m-2 p-2 flex flex-col align-center justify-items-center">
                <form className="flex flex-col gap-2 align-center justify-items-center" onSubmit={login}>
                    <label htmlFor="username">Username</label>
                    <input 
                        type="text"
                        name="username"
                        id="username"
                        className="rounded-md p-2"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}>
                    </input>

                    <button type="submit" className="bg-pink-200 rounded-md p-2">Login</button>
                </form>
            </div>
        </>
    )
}
