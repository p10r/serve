package de.p10r

import com.amazonaws.services.lambda.runtime.Context
import com.amazonaws.services.lambda.runtime.events.ScheduledEvent
import org.http4k.client.JavaHttpClient
import org.http4k.core.Response
import org.http4k.core.Uri
import org.http4k.core.then
import org.http4k.filter.ClientFilters
import org.http4k.serverless.AwsLambdaEventFunction
import org.http4k.serverless.FnHandler
import org.http4k.serverless.FnLoader


// The class name is referenced in `index.ts`
@Suppress("unused")
class ServeAppFunction : AwsLambdaEventFunction(EventFnLoader())

// The FnLoader is responsible for constructing the handler and for handling the serialisation of the request and response
fun EventFnLoader() = FnLoader { env: Map<String, String> ->
    FnHandler { _: ScheduledEvent, _: Context ->
        ProdApp(env)
    }
}

fun ProdApp(env: Map<String, String>): Response {
    val flashScoreUri = env["FLASH_SCORE_URI"]
        ?: error("env variable FLASH_SCORE_URI not provided")

    val discordUri = env["DISCORD_URI"]
        ?: error("env variable DISCORD_URI not provided")

    val apiKey = env["FLASH_SCORE_API_KEY"]
        ?: error("env variable DISCORD_URI not provided")

    return App(clientFrom(discordUri), clientFrom(flashScoreUri), apiKey)
}

private fun clientFrom(uri: String) = ClientFilters.SetBaseUriFrom(Uri.of(uri)).then(JavaHttpClient())
