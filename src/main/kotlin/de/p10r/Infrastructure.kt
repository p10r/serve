package de.p10r

import com.amazonaws.services.lambda.runtime.Context
import com.amazonaws.services.lambda.runtime.events.ScheduledEvent
import org.http4k.client.JavaHttpClient
import org.http4k.core.HttpHandler
import org.http4k.core.Uri
import org.http4k.core.then
import org.http4k.filter.ClientFilters
import org.http4k.filter.DebuggingFilters
import org.http4k.serverless.AwsLambdaEventFunction
import org.http4k.serverless.FnHandler
import org.http4k.serverless.FnLoader


//Http trigger
//@Suppress("unused")
//class ServeAppFunction : ApiGatewayV1LambdaFunction(AppLoader { env: Map<String, String> ->
//    ProdApp(env).routes()
//})

//Event trigger
@Suppress("unused")
class ServeAppFunction : AwsLambdaEventFunction(EventFnLoader(JavaHttpClient()))

// The FnLoader is responsible for constructing the handler and for handling the serialisation of the request and response
fun EventFnLoader(http: HttpHandler) = FnLoader { env: Map<String, String> ->
    println("in EventFnLoader")
    EventHandler(ProdApp(env))
}

fun EventHandler(app: App) = FnHandler { _: ScheduledEvent, _: Context ->
    println("starting workflow")
    app.run()
}


fun ProdApp(env: Map<String, String>): App {
    val flashScoreUri = env["FLASH_SCORE_URI"]
        ?: error("env variable FLASH_SCORE_URI not provided")

    val discordUri = env["DISCORD_URI"]
        ?: error("env variable DISCORD_URI not provided")

    val apiKey = env["FLASH_SCORE_API_KEY"]
        ?: error("env variable DISCORD_URI not provided")

    val flashScoreClient = ClientFilters.SetBaseUriFrom(Uri.of(flashScoreUri))
        .then(JavaHttpClient())

    val discordClient = ClientFilters.SetBaseUriFrom(Uri.of(discordUri))
        .then(DebuggingFilters.PrintRequestAndResponse())
        .then(JavaHttpClient())

    return App(discordClient, flashScoreClient, apiKey)
}