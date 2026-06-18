package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func yaziyiTemizle(hedefURL string) (readability.Article, error) {
	var bosYazi readability.Article

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", hedefURL, nil)
	if err != nil {
		return bosYazi, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return bosYazi, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return bosYazi, fmt.Errorf("site hata döndürdü: %d", resp.StatusCode)
	}

	parsedURL, _ := url.Parse(hedefURL)
	yazi, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		return bosYazi, err
	}

	return yazi, nil
}

const ortakCSSveHead = `
    <style>
        :root {
            --bg-body: #f4f6f8;
            --bg-container: #ffffff;
            --text-main: #2d3748;
            --text-muted: #718096;
            --border-color: #e2e8f0;
            --accent-color: #3182ce;
            --accent-hover: #2b6cb0;
            --code-bg: #edf2f7;
        }

        [data-theme="dark"] {
            --bg-body: #121212;
            --bg-container: #1e1e1e;
            --text-main: #e0e0e0;
            --text-muted: #888888;
            --border-color: #333333;
            --accent-color: #4dabf7;
            --accent-hover: #3793dd;
            --code-bg: #2d2d2d;
        }

        body {
            background-color: var(--bg-body);
            color: var(--text-main);
            transition: background-color 0.3s, color 0.3s;
        }
        
        .toggle-btn {
            background: var(--bg-container);
            border: 1px solid var(--border-color);
            color: var(--text-main);
            padding: 8px 15px;
            border-radius: 20px;
            cursor: pointer;
            font-size: 0.9rem;
            transition: all 0.3s;
        }
        .toggle-btn:hover {
            border-color: var(--accent-color);
        }
    </style>
    <script>
        const savedTheme = localStorage.getItem('theme') || 'system';
        if (savedTheme === 'dark' || (savedTheme === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
            document.documentElement.setAttribute('data-theme', 'dark');
        }
    </script>
`

const ortakJS = `
    <script>
        const i18n = {
            tr: {
                title: "Saf Okuyucu",
                subtitle: "İnternetin gürültüsünü susturun. Okumak istediğiniz yazının linkini yapıştırın.",
                placeholder: "https://ornek.com/yazi",
                readBtn: "Oku",
                back: "&larr; Ana Sayfaya Dön",
                sys: "💻 Sistem",
                dark: "🌙 Karanlık",
                light: "☀️ Aydınlık",
                langBtn: "🇹🇷 TR",
                docTitle: "Saf Okuyucu - Dijital Detoks"
            },
            en: {
                title: "Pure Reader",
                subtitle: "Silence the noise. Paste the link of the text you want to read.",
                placeholder: "https://example.com/text",
                readBtn: "Read",
                back: "&larr; Back to Home",
                sys: "💻 System",
                dark: "🌙 Dark",
                light: "☀️ Light",
                langBtn: "🇬🇧 EN",
                docTitle: "Pure Reader - Digital Detox"
            }
        };

        let currentLang = localStorage.getItem('lang') || 'tr';
        let currentTheme = localStorage.getItem('theme') || 'system';

        function updateUI() {
            document.documentElement.lang = currentLang;
            
            // Dil metinlerini güncelle
            document.querySelectorAll('[data-i18n]').forEach(el => {
                const key = el.getAttribute('data-i18n');
                if (el.tagName === 'INPUT') el.placeholder = i18n[currentLang][key];
                else el.innerHTML = i18n[currentLang][key];
            });

            // Buton içeriklerini güncelle
            const langBtn = document.getElementById('lang-btn');
            if(langBtn) langBtn.innerHTML = i18n[currentLang].langBtn;

            const themeBtn = document.getElementById('theme-btn');
            if(themeBtn) {
                if (currentTheme === 'dark') themeBtn.innerHTML = i18n[currentLang].dark;
                else if (currentTheme === 'light') themeBtn.innerHTML = i18n[currentLang].light;
                else themeBtn.innerHTML = i18n[currentLang].sys;
            }
        }

        function applyTheme(theme) {
            if (theme === 'dark') {
                document.documentElement.setAttribute('data-theme', 'dark');
            } else if (theme === 'light') {
                document.documentElement.removeAttribute('data-theme');
            } else {
                if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
                    document.documentElement.setAttribute('data-theme', 'dark');
                } else {
                    document.documentElement.removeAttribute('data-theme');
                }
            }
            updateUI();
        }

        document.addEventListener('DOMContentLoaded', () => {
            updateUI();
            
            const langBtn = document.getElementById('lang-btn');
            if (langBtn) {
                langBtn.addEventListener('click', () => {
                    currentLang = currentLang === 'tr' ? 'en' : 'tr';
                    localStorage.setItem('lang', currentLang);
                    updateUI();
                });
            }

            const themeBtn = document.getElementById('theme-btn');
            if (themeBtn) {
                themeBtn.addEventListener('click', () => {
                    if (currentTheme === 'system') currentTheme = 'dark';
                    else if (currentTheme === 'dark') currentTheme = 'light';
                    else currentTheme = 'system';
                    localStorage.setItem('theme', currentTheme);
                    applyTheme(currentTheme);
                });
            }
        });

        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
            if (currentTheme === 'system') applyTheme('system');
        });
    </script>
`

const htmlAnaSayfa = `
<!DOCTYPE html>
<html lang="tr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Saf Okuyucu - Dijital Detoks</title>
    ` + ortakCSSveHead + `
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; margin: 0; display: flex; align-items: center; justify-content: center; min-height: 100vh; }
        .top-bar { position: absolute; top: 20px; right: 20px; display: flex; gap: 10px; }
        .search-container { text-align: center; width: 100%; max-width: 600px; padding: 20px; }
        h1 { color: var(--text-main); font-size: 2.5rem; margin-bottom: 10px; font-weight: 400; }
        p { color: var(--text-muted); margin-bottom: 30px; }
        form { display: flex; gap: 10px; }
        input[type="url"] { flex: 1; background-color: var(--bg-container); border: 2px solid var(--border-color); border-radius: 8px; padding: 15px; color: var(--text-main); font-size: 1rem; outline: none; transition: border-color 0.3s; }
        input[type="url"]:focus { border-color: var(--accent-color); }
        button[type="submit"] { background-color: var(--accent-color); color: #fff; border: none; border-radius: 8px; padding: 0 25px; font-size: 1rem; font-weight: bold; cursor: pointer; transition: background-color 0.3s; }
        button[type="submit"]:hover { background-color: var(--accent-hover); }
    </style>
</head>
<body>
    <div class="top-bar">
        <button id="lang-btn" class="toggle-btn"></button>
        <button id="theme-btn" class="toggle-btn"></button>
    </div>
    <div class="search-container">
        <h1 data-i18n="title">Saf Okuyucu</h1>
        <p data-i18n="subtitle">İnternetin gürültüsünü susturun. Okumak istediğiniz yazının linkini yapıştırın.</p>
        <form action="/oku" method="GET">
            <input type="url" name="url" data-i18n="placeholder" placeholder="https://ornek.com/yazi" required autocomplete="off">
            <button type="submit" data-i18n="readBtn">Oku</button>
        </form>
    </div>
    ` + ortakJS + `
</body>
</html>
`

const htmlSablonu = `
<!DOCTYPE html>
<html lang="tr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    ` + ortakCSSveHead + `
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; line-height: 1.8; margin: 0; padding: 20px; }
        .container { max-width: 750px; margin: 40px auto; background: var(--bg-container); padding: 40px 50px; border-radius: 12px; box-shadow: 0 4px 15px rgba(0,0,0,0.1); transition: background-color 0.3s; }
        .header-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
        .top-bar-inline { display: flex; gap: 10px; }
        .header-bar a { color: var(--text-muted); font-size: 0.9rem; text-decoration: none; }
        .header-bar a:hover { color: var(--accent-color); }
        h1 { color: var(--text-main); border-bottom: 2px solid var(--border-color); padding-bottom: 15px; margin-bottom: 30px; font-size: 2rem; }
        h2, h3, h4 { color: var(--text-main); margin-top: 30px; }
        a { color: var(--accent-color); text-decoration: none; }
        a:hover { text-decoration: underline; }
        img { max-width: 100%; height: auto; border-radius: 8px; margin: 20px 0; display: block; }
        pre, code { background: var(--code-bg); padding: 4px 8px; border-radius: 6px; font-family: monospace; color: var(--text-main); }
        blockquote { border-left: 4px solid var(--accent-color); padding-left: 20px; color: var(--text-muted); font-style: italic; margin-left: 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header-bar">
            <a href="/" data-i18n="back">&larr; Ana Sayfaya Dön</a>
            <div class="top-bar-inline">
                <button id="lang-btn" class="toggle-btn"></button>
                <button id="theme-btn" class="toggle-btn"></button>
            </div>
        </div>
        <h1>{{.Title}}</h1>
        <div class="content">
            {{.Content}}
        </div>
    </div>
    ` + ortakJS + `
</body>
</html>
`

func main() {
	app := fiber.New()

	app.Use(limiter.New(limiter.Config{
		Max:        15,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).SendString("Çok fazla istek attınız. Lütfen bekleyin. / Too many requests. Please wait.")
		},
	}))

	tmpl := template.Must(template.New("okuyucu").Parse(htmlSablonu))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Type("html", "utf-8")
		return c.SendString(htmlAnaSayfa)
	})

	app.Get("/oku", func(c *fiber.Ctx) error {
		gelenURL := c.Query("url")
		if gelenURL == "" {
			return c.Redirect("/")
		}

		yazi, err := yaziyiTemizle(gelenURL)
		if err != nil {
			fmt.Printf("ASIL HATA BURADA -> %v\n", err)
			return c.Status(500).SendString("Hata oluştu: Yazı temizlenemedi. / Error: Text could not be parsed.")
		}

		data := struct {
			Title   string
			Content template.HTML
		}{
			Title:   yazi.Title,
			Content: template.HTML(yazi.Content),
		}

		c.Type("html", "utf-8")
		return tmpl.Execute(c.Response().BodyWriter(), data)
	})

	fmt.Println("Çift dilli ve temalı sürüm başlatıldı: http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
