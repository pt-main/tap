# Tap - библиотека для парсинга CLI на Rust

[![Crates.io](https://img.shields.io/crates/v/tap-rs.svg)](https://crates.io/crates/tap-rs)
[![Docs.rs](https://docs.rs/tap/badge.svg)](https://docs.rs/tap)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-yellow.svg)](https://opensource.org/licenses/Apache-2.0)

```bash
cargo install tap-rs
```

**Tap** - это лёгкая библиотека для построения CLI на Rust. Она предоставляет простой API для парсинга аргументов, работы с командами, флагами и цветным выводом, без внешних зависимостей и сложных макросов.

---

## Особенности

- **Цветной вывод** через короткие коды: `[?GN]`, `[?RD]`, `[?YW]` и т.д.
- **Автогенерация справки** с группировкой команд по одинаковым обработчикам.
- **Сабкоманды** (вложенные команды).
- **Нуль зависимостей** (только `std`).

---

## Быстрый старт

Создадим простое CLI с командой `hello`:

```rust
use tap::{
    argsparsing::ArgsP,
    engine::structs::{CmdInfo, Engine, ErrorType},
    utils::add_help_handler,
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

---

## Команды и аргументы

Добавьте команду с помощью `Engine::new_command`:

```rust
fn my_handler(e: &mut Engine, args: &[&str]) -> Result<(), Box<dyn Error>> { ... }

engine.new_command(
    "copy",                               // имя команды
    CmdInfo::new(
        my_handler,                       // функция-обработчик
        vec!["src", "dst"],               // обязательные аргументы
        vec!["force"],                    // опциональные
        false,                            // unlimited (если true, то принимает любое количество дополнительных аргументов)
        "Copy source to destination",     // описание для помощи
    ),
);
```

Помощь покажет:
```
╭───────  Command [copy]
⎬─ Args: <src>, <dst>, [force]
⎬─ Docs:
│     Copy source to destination
╰─────── 
```

---

## Флаги

Флаги записываются как `--flag`, `--key=value`, `--key:value`. Они парсятся автоматически и сохраняются в `ArgsP`:

- `args_p.flags` - список флагов без значений.
- `args_p.values` - словарь ключ -> значение.

В обработчике вы можете их прочитать:

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

---

## Цвета

Tap поддерживает цветной вывод через короткие коды вида `[?COLOR]`. Например:
- `[?GN]`, `[?GREEN]` - зелёный
- `[?BGN]`, `[?BGREEN]` - светлозеленый
- `[?BACKGREEN]`, `[?BKGN]` - зеленый фон
- `[?BACKBGREEN]`, `[?BKBGN]` - светлозеленый фон
- `[?BOLD]`, `[?BD]` - жирный
- `[?UNDERLINE]`, `[?UE]` - подчёркнутый
- `[?RT]` - сброс

Формат - `[?<BACKGROUND><BRIGHT><COLOR>]`, где:

- `<BACKGROUND>` - Для выбора цвета фона нужно добавить в начало цвета `BK` при сокращении цвета, `BACK` при полном названии.
- `<BRIGHT>` - Для использования светлой версии цвета нужно добавить в начало цвета `B` (сокращение от 'bright').
- `<COLOR>` - Варианты:
    - Две буквы: первая буква - начало названия цвета, вторая - конец названия цвета.
        Например: `GN`

Цвета - черный, красный, зеленый, желтый, синий, маджента, циан, и их светлые версии. Все цвета доступны для фона. Для текста доступен жирный стиль, подчеркнутый.

Используйте функции из `tap::formatting::print`, `tap::formatting::color`:

```rust
use tap::formatting::print::println;
use tap::formatting::color::colorize;

let raw = "[?GN]Success![?RT] File saved as [?YW]config.toml[?RT]";

println(raw); 
let colored = colorize(raw)
```

Цвета можно отключить глобально:

```rust
use tap::formatting::color::COLOR_ENABLED;
use std::sync::atomic::Ordering;

COLOR_ENABLED.store(false, Ordering::Relaxed);
```

---

## Группировка команд в помощи

Команды с одинаковыми значениями `CmdInfo` автоматически группируются в одну строку, как алиасы:

```rust
engine.new_commands(
    vec!["help", "-h"],
    CmdInfo::new(help_handler, vec![], vec![], false, "Show help"),
);
```

В помощи будет выведено: `[help / -h]`.


---

> MIT – подробности в файле LICENSE.
> 2026, By Pt.