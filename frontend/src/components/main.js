import React, { useState, useRef } from 'react';

function Main() {

    const [text, setText] = useState("")
    const [consola, setConsola] = useState("Consola: \n")

    const fileInputRef = useRef(null);

    const handleButtonAbrirClick = () => {
        fileInputRef.current.click();
    };

    const handleFileRead = (event) => {
        const content = event.target.result;
        setText(content);
    };

    const handleFileChosen = (file) => {
        const fileReader = new FileReader();
        fileReader.onloadend = handleFileRead;
        fileReader.readAsText(file);
    };

    const handleButtonEjecutarClick = async() => {
        const Comandos = text.split("\n")
        var Temp = "Consola: \n"
        for (let i = 0; i < Comandos.length; i++) {
            const Comando = Comandos[i]
            await fetch("http://localhost:4000/", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({string: Comando})
            })
                .then((response) => response.text())
                .then((data) => Temp += data)
                .catch((error) => console.log(error));
        }
        setConsola(Temp)
    };

    const handleTextChange = (event) => {
        setText(event.target.value);
    };

    return (
        <div>
            <div id="Buttons" className="base">
                <form className="main_form">
                    <button className="col-4 tecla" type="button" onClick={handleButtonAbrirClick}>Abrir</button>
                    <input
                        type="file"
                        ref={fileInputRef}
                        onChange={(e) => handleFileChosen(e.target.files[0])}
                        style={{ display: 'none' }}
                    />
                    <button className="col-4 tecla" type="button" onClick={handleButtonEjecutarClick}>Ejecutar</button>
                    <button className="col-4 tecla" type="button">Login</button>
                </form>
            </div>
            <div id="Texts" className="base">
                <form className="main_form">
                    <textarea className="text" type="button" value={text} onChange={handleTextChange} />
                    <textarea readOnly={true} className="text" type="button" value={consola} />
                </form>
            </div>
        </div>
    );
}

export default Main;