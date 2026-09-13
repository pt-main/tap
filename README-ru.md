# Tap - библиотека для парсинга CLI на Rust

[![Crates.io](https://img.shields.io/crates/v/tap-rs.svg)](https://crates.io/crates/tap-rs)
[![GitHub](https://img.shields.io/badge/GitHub-repo-181717?logo=github)](https://github.com/pt-main/tap/tree/rust)
[![Docs.rs](https://docs.rs/tap-rs/badge.svg)](https://docs.rs/tap-rs)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-yellow.svg)](https://opensource.org/licenses/Apache-2.0)

```bash
cargo install tap-rs
cargo add tap-rs
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
fn my_handler(e: &mut Engine, args: &[&str]) -> ErrorType { ... }

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

Tap поддерживает цветной вывод с помощью коротких кодов вроде `[?COLOR]`. Примеры:
- `[?GN]`, `[?GREEN]` - зелёный
- `[?BGN]`, `[?BGREEN]` - ярко-зелёный
- `[?BACKGREEN]`, `[?BKGN]` - зелёный фон
- `[?BACKBGREEN]`, `[?BKBGN]` - ярко-зелёный фон
- `[?BOLD]`, `[?BD]` - жирный
- `[?DIM]`, `[?DM]` - тусклый / приглушённый
- `[?ITALIC]`, `[?IT]` - курсив
- `[?UNDERLINE]`, `[?UE]` - подчёркнутый
- `[?BLINK]`, `[?BL]` - мигающий
- `[?REVERSE]`, `[?RV]` - инвертированный
- `[?STRIKETHROUGH]`, `[?ST]` - зачёркнутый
- `[?RT]` - сброс

Формат - `[?<BACKGROUND><BRIGHT><COLOR>]`, где:

- `<BACKGROUND>` - чтобы выбрать цвет фона, добавьте перед цветом `BK` (кратко) или `BACK` (полное имя).
- `<BRIGHT>` - чтобы использовать яркий вариант, добавьте `B` (сокращение от 'bright').
- `<COLOR>` - варианты:
    - Две буквы: первая буква - начальная буква названия цвета, вторая - конечная.
        Пример: `GN`

### Цвета

**Стандартные 16** - чёрный, красный, зелёный, жёлтый, синий, пурпурный, голубой, белый и их яркие варианты. Все цвета доступны для фона.

**Расширенные (xterm-256)** - дополнительные именованные цвета, доступные как для переднего плана, так и для фона:

| Название     | Сокр. | Название     | Сокр. |
|--------------|-------|--------------|-------|
| оранжевый    | `OE`  | лаймовый     | `LE`  |
| розовый      | `PK`  | бирюзовый    | `TL`  |
| пурпурный    | `PE`  | тёмно-синий  | `NY`  |
| фиолетовый   | `VT`  | золотой      | `GD`  |
| коричневый   | `BN`  | серебряный   | `SR`  |
| серый        | `GY`  | бордовый     | `MN`  |

Примеры: `[?OE]` - оранжевый, `[?BACKGOLD]` / `[?BKGD]` - золотой фон, `[?BKOE]` - оранжевый фон.

### Стили

Для текста доступны следующие стили: `BOLD` (`BD`), `DIM`/`FAINT` (`DM`), `ITALIC` (`IT`), `UNDERLINE` (`UE`), `BLINK` (`BL`), `REVERSE` (`RV`), `STRIKETHROUGH` (`ST`) и `RESET` (`RT`).

### Использование

Используйте функции из `tap::formatting::print` и `tap::formatting::color`:

```rust
use tap::formatting::print::println;
use tap::formatting::color::colorize;

let raw = "[?GN]Успех![?RT] Файл сохранён как [?YW]config.toml[?RT]";

println(raw); 
let colored = colorize(raw)
```

Цвета можно отключить глобально:

```rust
use tap::formatting::color::COLOR_ENABLED;
use std::sync::atomic::Ordering;

COLOR_ENABLED.store(false, Ordering::Relaxed);
```

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