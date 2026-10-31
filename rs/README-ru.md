# Tap - библиотека для парсинга CLI на Rust

[![Crates.io](https://img.shields.io/crates/v/tap-rs.svg)](https://crates.io/crates/tap-rs)
[![GitHub](https://img.shields.io/badge/GitHub-repo-181717?logo=github)](https://github.com/pt-main/tap/tree/rust)
[![Docs.rs](https://docs.rs/tap-rs/badge.svg)](https://docs.rs/tap-rs)

```bash
cargo add tap-rs
```

**Tap** - небольшая библиотека для построения CLI на Rust. Парсит аргументы, команды и флаги, печатает цветной текст, генерирует справку. Одна зависимость `string_template`, макросы опциональны.

## Возможности

- Цветной вывод через короткие коды `[?GN]`, `[?RD]`, `[?YW]`.
- Автогенерация справки с группировкой команд по одинаковым обработчикам.
- Сабкоманды.
- Типизированный разбор позиционных аргументов через `ArgsExt`.
- Флаги `--flag`, `--key=value`, `--key:value` в `ArgsP`.
- Макросы `commands!`, `subcommands!`, `args!`, `bail!`, `ensure!`, `cprintln!` и другие.
- Рендереры с шаблонами, типизацией и JSON-сериализацией.

## Быстрый старт

Простое CLI с командой `hello`:

```rust
use tap::{
    argsparsing::ArgsP,
    engine::structs::{CmdInfo, Engine, ErrorType},
    tooling::add_help_handler,
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

Символы рамки в выводе `help` - это то, что реально печатает `generate_help`. Это не украшение README.

## Быстрый старт с макросами

То же самое, но короче:

```rust
use tap::{
    argsparsing::ArgsP, commands,
    engine::structs::{Engine, ErrorType},
    tooling::add_help_handler,
};

fn hello_handler(_: &mut Engine, args: &[&str]) -> ErrorType {
    let name = args.first().copied().unwrap_or("World");
    println!("Hello, {}!", name);
    Ok(())
}

fn main() {
    let mut ap = ArgsP::new();
    let mut engine = Engine::new(&mut ap, "MyApp v1.0".to_string());

    commands! {
        engine;
        "hello" => hello_handler { optional: ["name"], docs: "Say hello" }
    }

    add_help_handler(vec!["help", "-h"], &mut engine);

    if let Err(e) = engine.main() {
        eprintln!("Error: {}", e);
    }
}
```

Макрос `commands!` сам создаёт `CmdInfo` и регистрирует команду. Смотрите раздел Макросы ниже.

## Команды и аргументы

```rust
fn my_handler(e: &mut Engine, args: &[&str]) -> ErrorType { ... }

engine.new_command(
    "copy",
    CmdInfo::new(
        my_handler,                    // обработчик
        vec!["src", "dst"],            // обязательные
        vec!["force"],                 // опциональные
        false,                         // unlimited: true - принимает любое число аргументов
        "Copy source to destination",  // описание для справки
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

### Сабкоманды

Сабкоманды - это команды, которые получают только аргументы, без `&mut Engine`. Их обработчик имеет тип `fn(&[&str]) -> ErrorType`.

```rust
fn echo_handler(raw: &[&str]) -> ErrorType {
    println!("echo: {}", raw.join(" "));
    Ok(())
}

engine.new_subcommand(
    "echo",
    CmdInfo::new(echo_handler, vec![], vec![], true, "Echo back args"),
);
```

Регистрация через `new_subcommand` / `new_subcommands` или макрос `subcommands!`.

### Разбор аргументов через `ArgsExt`

Трейт `ArgsExt` даёт типизированный доступ к срезу аргументов:

```rust
use tap::argsparsing::ArgsExt;

fn add_handler(_: &mut Engine, raw: &[&str]) -> ErrorType {
    let a: f64 = raw.arg(0)?;
    let b: f64 = raw.arg(1)?;
    let label = raw.arg_or(2, "sum");   // если аргумента нет - вернёт "sum"
    println!("{}: {} + {} = {}", label, a, b, a + b);
    Ok(())
}
```

Ошибки содержат индекс аргумента, а не голую ошибку парсера.

## Флаги

Флаги записываются как `--flag`, `--key=value` или `--key:value`. Разбираются автоматически в `ArgsP`:

- `args_p.flags` - список флагов без значений.
- `args_p.values` - словарь ключ -> значение.

Флаги вырезаются из `args` до вызова обработчика. Смотрите их только через `e.args_p` или вспомогательные методы `Engine`.

```rust
fn my_handler(e: &mut Engine, _: &[&str]) -> ErrorType {
    if e.flag("verbose") {
        println!("verbose mode");
    }
    if let Some(out) = e.value("output") {
        println!("output file: {}", out);
    }
    let workers: u16 = e.value_t("workers")?.unwrap_or(4);
    println!("workers: {}", workers);
    Ok(())
}
```

- `Engine::flag(name)` - `bool`.
- `Engine::value(name)` - `Option<&str>`.
- `Engine::value_t::<T>(name)` - `Result<Option<T>, String>`; ошибка парсинга приходит с именем флага.

### `--` как разделитель

`--` отключает разбор флагов. Всё после него попадает в `args` как есть:

```
myapp run -- --x --y    # args = ["--x", "--y"], flags = []
```

## Макросы

Все макросы экспортируются из корня крейта, то есть вызываются как `tap::commands!`, `tap::args!` и т.д.

### `commands!` и `subcommands!`

```rust
tap::commands! {
    engine;
    "greet" => greet_handler {
        optional: ["name"],
        docs: "Greet someone",
    }
    "add" | "sum" => add_handler {
        required: ["a", "b"],
        docs: "Sum two numbers",
    }
    "echo" => echo_handler {
        unlimited: true,
        docs: "Echo args",
    }
}
```

Поля `required`, `optional`, `unlimited` опциональны, `docs` обязателен. `"name" | "alias"` регистрирует несколько имён с одним обработчиком.

`subcommands!` работает так же, но обработчик должен быть `fn(&[&str]) -> ErrorType`.

### `args!`

Привязывает позиционные аргументы к типизированным переменным:

```rust
fn add_handler(_: &mut Engine, raw: &[&str]) -> ErrorType {
    tap::args! { raw;
        a: f64,
        b: f64,
        label: String = "sum".to_string(),
    }
    println!("{}: {} + {} = {}", label, a, b, a + b);
    Ok(())
}
```

Форма `name: Type` - обязательный аргумент, `name: Type = default` - опциональный. На каждый аргумент расходуется одна позиция, порядок важен. Ошибка приходит с именем аргумента и индексом.

### `bail!` и `ensure!`

```rust
fn h(_: &mut Engine, raw: &[&str]) -> ErrorType {
    tap::ensure!(!raw.is_empty(), "at least one argument required");
    let n: u32 = raw[0].parse().map_err(|_| format!("bad number: {}", raw[0]))?;
    tap::ensure!(n > 0, "must be positive, got {}", n);
    Ok(())
}
```

`bail!` досрочно возвращает `Err` из обработчика. `ensure!` - то же самое, но по условию.

### `cprint!`, `cprintln!`, `ceprintln!`

Аналоги `print!` / `println!` / `eprintln!`, но с автоматической раскраской `[?COLOR]`-кодов:

```rust
tap::cprintln!("[?GN]Success![?RT] File saved as [?YW]config.toml[?RT]");
```

### `scope!`, `template!`, `template_colored!`

Сокращения для работы с `Renderable`:

```rust
use tap::tooling::renderable::Renderable;

let mut r = Renderable::new();
tap::scope!(r;
    "name"        => "Rust",
    "version:int" => "1",
    "stable:bool" => "true",
);
tap::template!(r, "plain", "Hello, {{name}} (v{{version}})");
tap::template_colored!(r, "colored", "[?GN]Hello, {{name}}[?RT]");
```

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

- `<BACKGROUND>` - префикс `BK` или `BACK` для фона, опционально.
- `<BRIGHT>` - префикс `B` для яркого варианта, опционально.
- `<COLOR>` - имя или две буквы: первая и последняя буква названия, написанные заглавными буквами. `GN` = green.

### Палитра

**Стандартные 16**: black, red, green, yellow, blue, magenta, cyan, white и их яркие варианты. Все доступны как фон.

**Расширенные xterm-256** (доступны и как фон):

| Название | Сокр. | Название | Сокр. |
|----------|-------|----------|-------|
| orange   | `OE`  | lime     | `LE`  |
| pink     | `PK`  | teal     | `TL`  |
| purple   | `PE`  | navy     | `NY`  |
| violet   | `VT`  | gold     | `GD`  |
| brown    | `BN`  | silver   | `SR`  |
| gray     | `GY`  | maroon   | `MN`  |

Примеры: `[?OE]` - оранжевый, `[?BACKGOLD]` / `[?BKGD]` - золотой фон, `[?BKOE]` - оранжевый фон.

### Использование

```rust
use tap::tooling::colors::print::println;
use tap::tooling::colors::set::colorize;

let raw = "[?GN]Успех![?RT] Файл сохранён как [?YW]config.toml[?RT]";

println(raw);                        // печатает уже раскрашенным
let colored: String = colorize(raw); // или получите строку
```

`println` и `print` в `colors::print` уже вызывают `colorize` внутри - не нужно красить дважды.

Отключить цвета глобально:

```rust
use tap::tooling::colors::set::COLOR_ENABLED;
use std::sync::atomic::Ordering;

COLOR_ENABLED.store(false, Ordering::Relaxed);
```

Неизвестные коды `[?FOO]` не съедаются - остаются в строке как есть.

## Группировка команд в справке

Команды с одинаковыми `CmdInfo` автоматически группируются в одну строку как алиасы:

```rust
engine.new_commands(
    vec!["help", "-h"],
    CmdInfo::new(help_handler, vec![], vec![], false, "Show help"),
);
```

В справке будет `[help / -h]`.

Сравниваются `required`, `optional`, `unlimited` и `docs`. Если что либо из них отличается, группировки не будет.

## Рендереры

Модуль `tooling::renderable` шаблонизирует вывод по scope:

```rust
use tap::tooling::renderable::{
    Renderable, add_template_render, add_template_render_colored,
};

let mut r = Renderable::new();
r.scope.insert("name", "Rust");
r.scope.insert("version:int", "1");
r.scope.insert("stable:bool", "true");

add_template_render("Hello, {{name}} (v{{version}})", "plain", &mut r);
add_template_render_colored("[?GN]Hello, {{name}}[?RT]", "colored", &mut r);

println!("{}", r.render_with("plain").unwrap());
println!("{}", r.render_with("colored").unwrap());
println!("{}", r.render_json().unwrap());   // {"name":"Rust","stable":true,"version":1}
```

Суффикс ключа задаёт тип при сериализации в JSON:

- `:int` - целое; строка должна совпадать с `-?(0|[1-9][0-9]*)`.
- `:float` - число с плавающей точкой и опциональной экспонентой.
- `:bool` - только `"true"` или `"false"`.
- `:raw` - значение вставляется в JSON как есть, без кавычек и проверки.
- Без суффикса - строка.

`render_json` возвращает `Result`. При невалидном значении вернёт `RenderError::InvalidValue`, а не молча вставит мусор в JSON.

Для массива строк есть `render_json_array(&[Renderable])`:

```rust
use tap::tooling::renderable::render_json_array;

let json = render_json_array(&rows)?;   // [{"...":...},{"...":...}]
```

## Пример

В репозитории есть `examples/demo.rs`. Он показывает все основные возможности: макросы, флаги, типизированные аргументы, сабкоманды, цвета, рендереры и JSON.

```bash
cargo run --example demo -- help
cargo run --example demo -- greet Rust
cargo run --example demo -- greet --loud Rust
cargo run --example demo -- add 2 3.5
cargo run --example demo -- serve 127.0.0.1 8080 --verbose --tls=on --workers=8
cargo run --example demo -- report
cargo run --example demo -- echo hello world
```

---

> MIT - подробности в файле LICENSE.
> 2026, By Pt.