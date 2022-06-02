package de.p10r

import java.time.Instant
import java.time.LocalDate
import java.time.LocalTime
import java.time.ZoneId
import java.time.format.DateTimeFormatter
import java.time.format.FormatStyle

fun discordMessageOf(
    schedules: Schedules,
    today: LocalDate,
): DiscordMessage {
    val fields = schedules.map { league ->
        DiscordMessage.Embeds.Fields(
            name = league.leagueName.localized(),
            value = league.games.toText(),
            inline = false
        )
    }

    return DiscordMessage(
        content = "Games for ${today.format()}",
        embeds = listOf(DiscordMessage.Embeds(fields))
    )
}

private fun List<Game>.toText() = joinToString("\n") {
    "**${it.homeTeam} - ${it.awayTeam}** \t ${it.time.localize()}"
}

private fun Instant.localize() = mapOf(
    "Europe/Berlin" to "BER",
    "America/New_York" to "NY",
    "America/Los_Angeles" to "LA",
    "Asia/Hong_Kong" to "HK",
).map { (zoneId, name) ->
    this.atZone(ZoneId.of(zoneId))
        .toLocalTime()
        .format() + " $name"
}.joinToString(prefix = "(", separator = "/", postfix = ")")

private fun LocalDate.format() = format(DateTimeFormatter.ofLocalizedDate(FormatStyle.FULL))

private fun LocalTime.format() = format(DateTimeFormatter.ofPattern("HH:mm"))

private fun String.localized() = when {
    contains("Poland")  -> "🇵🇱 $this"
    contains("Italy")   -> "🇮🇹 $this"
    contains("France")  -> "🇫🇷 $this"
    contains("Germany") -> "🇩🇪 $this"
    contains("Russia")  -> "🇷🇺 $this"
    else                -> this
}
