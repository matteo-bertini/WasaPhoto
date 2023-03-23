<script>
export default {
	// Dichiarazione del reactive state
	data: function () {
		return {
			ErrorMessage: null,
			Username : "",
			
		}
	},
	methods: {
		// Click sul pulsante Login
		async LoginButtonPressed(){
			try{
				
				// DoLogin
				let doLogin_response = await this.$axios.post("/session",{Username:this.Username});
				localStorage.setItem("Authstring",doLogin_response.data.Identifier);
				localStorage.setItem("Username",this.Username);
				

				// GetUserProfile per vedere se devo o meno creare il profilo
				let getUserProfile_config = {
					headers: {
						Authorization: `Bearer ${localStorage.getItem("Authstring")}`
					},
					params: {
						Username: this.Username
					},
					validateStatus: function (status) {
						 return status <500; 
  					},
				};
				let getUserProfile_response =  await this.$axios.get("/users/",getUserProfile_config);
				
				// Se la risposta è 404,devo creare il profilo
				if(getUserProfile_response.status==404){
					
					// AddUser
					let addUser_config = {
						headers: {
							Authorization: `Bearer ${localStorage.getItem("Authstring")}`
						}
					}
					try{
						let addUser_response = await this.$axios.post("/users/",{Username:this.Username},addUser_config);
						this.$router.push("/users/"+this.Username+"/");
						return;
					}
					catch(e){
						console.log(e);
						return;
					}


				}
				this.$router.push("/users/"+this.Username+"/");
				return;


			}catch(e){
				console.log(e);
				return;
			}


		}
	}
	
}


</script>

<template>
	<div class="container" style="display: flex; flex-direction: column; width:100vw; height: 100vh; justify-content: flex-start; align-items: center;">
		
		<!-- Titolo -->
		<div>
			<span class="badge badge-secondary" style="display: flex; flex-direction: column; justify-content: center; align-items: center; width: 40vw; height: 20vh; font-size: 5em; color: rgb(0, 0, 0); background-color: white;">
				WasaPhoto
			
			</span>
		</div>
		
		
		
		<!-- Pannello di Login -->
		<div class="card" style="display: flex; flex-direction: column; position: absolute; top: 35%; width: 50vw; height: 20%; justify-content: center; align-items: center; background-color:rgba(255, 255, 255, 0.285);">
			<span style="font-family:Arial; font-size:xx-large; font-weight: bolder;">LOGIN</span>
			<input id="LoginTextArea" type="text" class="form-control" v-model="Username" style="width: 40%;" placeholder="Type here your username...">
			<button @click="LoginButtonPressed" id="LoginButton" class="btn btn-dark" style="margin-top: 2em;" >
				<i class="fa-solid fa-right-to-bracket"> Login/Register </i>
			
			</button>	
		</div>

		



		
	</div>
	
</template>

<style>
	#LoginButton:hover{
		transform: scale(1.1);

	}
	

	
</style>
