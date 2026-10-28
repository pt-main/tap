# Tap - библиотека для парсинга CLI на Rust

[![Crates.io](https://img.shields.io/crates/v/tap-rs.svg)](https://crates.io/crates/tap-rs)
[![GitHub](https://img.shields.io/badge/GitHub-repo-181717?logo=github)](https://github.com/pt-main/tap/tree/rust)
[![Docs.rs](https://docs.rs/tap-rs/badge.svg)](https://docs.rs/tap-rs)

```bash
cargo add tap-rs
# OR
cargo add --git https://github.com/pt-main/tap/rs
```

**Tap** - лёгкая библиотека для построения CLI на Rust. Простой API для парсинга аргументов, команд, флагов и цветного вывода. Одна зависимость (`string_template`), без макросов.

---

## Особенности

- Цветной вывод через короткие коды: `[?GN]`, `[?RD]`, `[?YW]` и т.д.
- Автогенерация справки с группировкой команд по одинаковым обработчикам.
- Сабкоманды (вложенные команды).
- Рендереры вывода с шаблонами и типизацией плюс JSON-сериализация.

---

## Быстрый старт

Простое CLI с командой `hello`:

```rust
use tap::{
    argsparsing::ArgsP,
    engine::structs::{CmdInfo, Engine, ErrorType},
    tooling::console::engine::add_help_handler,
};

fn hello_handler(_: &mut Engine, args: &[&str]) -> ErrorType {
    let name = args.first().unwrap_or(&"World");
    println!("Hello, {}!", name);
    Ok(())
}

fn main() {
    let mut ap = ArgsP::new();
    let mut engine = Engine::new(&mut ap, "MyApp v1.0".to_string());

    engine.new_command(
        "hello",
        CmdInfo::new(hello_handler, vec![], vec!["name"], false, "Say hello"),
    );

    add_help_handler(vec!["help", "-h"], &mut engine);

    if let Err(e) = engine.main() {
        eprintln!("Error: {}", e);
    }
}
```

Запуск:

```bash
$ cargo run hello Rust
Hello, Rust!
$ cargo run hello
Hello, World!
$ cargo run help
MyApp v1.0

╭───────  Command [help OR -h]
⎬─ Args: (none)
⎬─ Docs:
│     Generate and print help message
╰─────── 
╭───────  Command [hello]
⎬─ Args: [name]
⎬─ Docs:
│     Say hello
╰─────── 
```

Символы рамки в выводе `help` - это то, что реально печатает `generate_help`. Не украшение README.

---

## Команды и аргументы

```rust
fn my_handler(e: &mut Engine, args: &[&str]) -> ErrorType { ... }

engine.new_command(
    "copy",                               // имя команды
    CmdInfo::new(
        my_handler,                       // функция-обработчик
        vec!["src", "dst"],               // обязательные аргументы
        vec!["force"],                    // опциональные
        false,                            // unlimited: если true, принимает любое число доп. аргументов
        "Copy source to destination",     // описание для справки
    ),
);
```

Справка покажет:

```
╭───────  Command [copy]
⎬─ Args: <src>, <dst>, [force]
⎬─ Docs:
│     Copy source to destination
╰─────── 
```

---

## Флаги

Флаги записываются как `--flag`, `--key=value`, `--key:value`. Парсятся автоматически в `ArgsP`:

- `args_p.flags` - список флагов без значений.
- `args_p.values` - словарь ключ -> значение.

Чтение в обработчике:

```rust
fn my_handler(e: &mut Engine, args: &[&str]) -> ErrorType {
    if e.args_p.flags.contains(&"verbose".to_string()) {
        println!("Verbose mode enabled");
    }
    if let Some(out) = e.args_p.values.get("output") {
        println!("Output file: {}", out);
    }
    Ok(())
}
```

Важно: флаги вырезаются из `args` до вызова обработчика. Смотрите их только через `e.args_p`.

---

## Цвета

Цветной вывод через коды `[?COLOR]`:

- `[?GN]`, `[?GREEN]` - зелёный
- `[?BGN]`, `[?BGREEN]` - ярко-зелёный
- `[?BACKGREEN]`, `[?BKGN]` - зелёный фон
- `[?BACKBGREEN]`, `[?BKBGN]` - ярко-зелёный фон
- `[?BOLD]`, `[?BD]` - жирный
- `[?DIM]`, `[?DM]` - тусклый
- `[?ITALIC]`, `[?IT]` - курсив
- `[?UNDERLINE]`, `[?UE]` - подчёркнутый
- `[?BLINK]`, `[?BL]` - мигающий
- `[?REVERSE]`, `[?RV]` - инвертированный
- `[?STRIKETHROUGH]`, `[?ST]` - зачёркнутый
- `[?RT]` - сброс

Формат: `[?<BACKGROUND><BRIGHT><COLOR>]`.

- `<BACKGROUND>` - префикс `BK` (кратко) или `BACK` (полностью) для фона. (опциональный)
- `<BRIGHT>` - префикс `B` для яркого варианта. (опциональный)
- `<COLOR>` - имя или две буквы (первая и последняя буква названия). Пример: `GN` = green. (обязательный)

### Палитра

**Стандартные 16**: black, red, green, yellow, blue, magenta, cyan, white и их яркие варианты. Все доступны как фон.

**Расширенные (xterm-256)**: доступны и как фон:

| Название  | Сокр. | Название    | Сокр. |
|-----------|-------|-------------|-------|
| orange    | `OE`  | lime        | `LE`  |
| pink      | `PK`  | teal        | `TL`  |
| purple    | `PE`  | navy        | `NY`  |
| violet    | `VT`  | gold        | `GD`  |
| brown     | `BN`  | silver      | `SR`  |
| gray      | `GY`  | maroon      | `MN`  |

Примеры: `[?OE]` - оранжевый, `[?BACKGOLD]` / `[?BKGD]` - золотой фон, `[?BKOE]` - оранжевый фон.

### Стили

`BOLD` (`BD`), `DIM`/`FAINT` (`DM`), `ITALIC` (`IT`), `UNDERLINE` (`UE`), `BLINK` (`BL`), `REVERSE` (`RV`), `STRIKETHROUGH` (`ST`), `RESET` (`RT`).

### Использование

```rust
use tap::tooling::colors::print::println;
use tap::tooling::colors::set::colorize;

let raw = "[?GN]Успех![?RT] Файл сохранён как [?YW]config.toml[?RT]";

println(raw);                       // печатает уже раскрашенным
let colored: String = colorize(raw); // или получите строку
```

`println`/`print` в `colors::print` уже вызывают `colorize` внутри - не нужно красить дважды.

Отключить цвета глобально:

```rust
use tap::tooling::colors::set::COLOR_ENABLED;
use std::sync::atomic::Ordering;

COLOR_ENABLED.store(false, Ordering::Relaxed);
```

Неизвестные коды (`[?FOO]`) не съедаются и не превращаются в мусор - они остаются в строке как есть. Это удобно для отладки.

---

## Группировка команд в справке

Команды с одинаковыми `CmdInfo` автоматически группируются в одну строку как алиасы:

```rust
engine.new_commands(
    vec!["help", "-h"],
    CmdInfo::new(help_handler, vec![], vec![], false, "Show help"),
);
```

В справке будет `[help / -h]`.

Сравниваются `required`, `optional`, `unlimited` и `docs`. Если `docs` отличаются - группировки не будет.

---

## Рендереры

Модуль `tooling::renderable` позволяет шаблонизировать вывод по scope:

```rust
use tap::tooling::renderable::{
    Renderable, add_template_render, add_template_render_colored,
};

let mut r = Renderable::new();
r.scope.insert("name", "Rust");
r.scope.insert("version:int", "1");     // :int - валидируется как число
r.scope.insert("stable:bool", "true");  // :bool - только "true"/"false"

add_template_render("Hello, {{name}} (v{{version}})", "plain", &mut r);
add_template_render_colored("[?GN]Hello, {{name}}[?RT]", "colored", &mut r);

println!("{}", r.render_with("plain").unwrap());
println!("{}", r.render_with("colored").unwrap());
println!("{}", r.render_json().unwrap()); // {"name":"Rust","stable":true,"version":1}
```

`render_json` возвращает `Result` - при невалидном `:int`/`:bool` вернёт `RenderError::InvalidValue`, а не молча вставит мусор в JSON.

---

> MIT - подробности в файле LICENSE.
> 2026, By Pt.