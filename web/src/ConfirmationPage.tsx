import {API_URL} from "./App.tsx";
import {useNavigate, useParams} from "react-router-dom";

export const ConfirmationPageComponent = () => {
    const { token } = useParams()
    const redirect = useNavigate()
    const handleConfirm =  async () => {
        const response = await fetch(`${API_URL}/users/activate/${token}`, {
            method: "PUT"
        })
        if (response.ok) {
            // redirect to the home page
            redirect("/")
        } else {
            alert("Failed to confirm token")
        }
    }

    return (
        <div>
            <h2>Confirmation</h2>
            <button onClick={handleConfirm}>Click to confirm</button>
        </div>
    )
}
