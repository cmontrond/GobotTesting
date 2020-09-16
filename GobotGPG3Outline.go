package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	g "gobot.io/x/gobot/platforms/dexter/gopigo3"
	"gobot.io/x/gobot/platforms/raspi"
	"time"
)

//robotRunLoop is the main function for the robot, the gobot framework
//will spawn a new thread in the NewRobot factory functin and run this
//function in that new thread. Do all of your work in this function and
//in other functions that this function calls. don't read from sensors or
//use actuators frmo main or you will get a panic.
func robotRunLoop(gopigo3 *g.Driver, lightSensor *aio.GroveLightSensorDriver, soundSensor *aio.GroveSoundSensorDriver) {
	for {

		soundSensorVal, soundSensorErr := soundSensor.Read()

		if soundSensorErr != nil {
			fmt.Errorf("Error reading sensor %+v", soundSensorErr)
		}
		fmt.Println("Sound Value is ", soundSensorVal)

		//sensorVal, err := lightSensor.Read()
		//if err != nil {
		//	fmt.Errorf("Error reading sensor %+v", err)
		//}
		//fmt.Println("Light Value is ", sensorVal)

		//if sensorVal > 1000 {
		//	//_ = gopigo3.SetLED(g.LED_EYE_RIGHT, 0x00, 0x00, uint8(0xFF))
		//	result, err := soundSensor.Read()
		//} else {
		//	//_ = gopigo3.SetLED(g.LED_EYE_RIGHT, 0x00, 0x00, uint8(0x00))
		//}

		time.Sleep(time.Second)
	}
}

func main() {
	//We create the adaptors to connect the GoPiGo3 board with the Raspberry Pi 3
	//also create any sensor drivers here
	raspiAdaptor := raspi.NewAdaptor()
	gopigo3 := g.NewDriver(raspiAdaptor)
	lightSensor := aio.NewGroveLightSensorDriver(gopigo3, "AD_2_1") //AnalogDigital Port 1 is "AD_1_1" this is port 2
	soundSensor := aio.NewGroveSoundSensorDriver(gopigo3, "AD_1_1")
	//end create hardware drivers

	//here we create an anonymous function assigned to a local variable
	//the robot framework will create a new thread and run this function
	//I'm calling my robot main loop here. Pass any of the variables we created
	//above to that function if you need them
	mainRobotFunc := func() {
		robotRunLoop(gopigo3, lightSensor, soundSensor)
	}

	//this is the crux of the gobot framework. The factory function to create a new robot
	//struct (go uses structs and not objects) It takes four parameters
	robot := gobot.NewRobot("gopigo3sensorChecker", //first a name
		[]gobot.Connection{raspiAdaptor},     //next a slice of connections to one or more robot controllers
		[]gobot.Device{gopigo3, lightSensor}, //next a slice of one or more sensors and actuators for the robots
		mainRobotFunc,                        //the variable holding the function to run in a new thread as the main function
	)

	robot.Start() //actually run the function
}
