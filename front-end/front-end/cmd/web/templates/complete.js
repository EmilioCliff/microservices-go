let feedback = document.getElementById("payment-message");

fetch("http:\/\/localhost:8080/handler", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
        action: "payment",
        payment: {
            publish: "request publishable-key",
        }
    }),
})
.then(response => response.json())
.then((data) => {
    stripe = Stripe(data.publishable_key)
    checkStatusHere();
})
.catch((error) => {
    feedback.innerHTML =  "<br><br>Error:" + error;
})

async function checkStatusHere() {
    const clientSecret = new URLSearchParams(window.location.search).get(
        "payment_intent_client_secret"
    );

    const { paymentIntent } = await stripe.retrievePaymentIntent(clientSecret);
    
    let intentP = document.querySelector(".payment-intent");
    intentP.innerText = JSON.stringify(paymentIntent, null, 2);
    
    switch (paymentIntent.status) {
        case "succeeded":
          feedback.innerText = "Payment succeeded!"
          break;
        case "processing":
          feedback.innerText = "Your payment is processing."
          break;
        case "requires_payment_method":
          feedback.innerText = "Your payment was not successful, please try again."
          break;
        default:
          feedback.innerText = "Something went wrong."
          break;
    }
}