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
	<div class="container-fluid" id="LoginPageContainer" style="display: flex; flex-direction: column; width:100vw; height: 100vh; justify-content: flex-start; align-items: center;">
		
		<!-- Titolo -->
		<div style="margin-top: 5rem; font-size: 2rem;">
			<i class="fa-solid fa-camera"> WasaPhoto </i>
		</div>
		
		
		
		<!-- Pannello di Login -->
		<div class="card" style="display: flex; flex-direction: column; position: absolute; top: 35%; width: 50vw; height: 20%; justify-content: center; align-items: center; background-color:rgba(255, 255, 255, 0.285);">
			<span  id="LoginPageTitleSpan" style="font-family:Helvetica; font-size:xx-large; font-weight: bolder;">Login</span>
			<input id="LoginTextArea" type="text" class="form-control" v-model="Username" style="width: 40%;" placeholder="Type here your username...">
			<button @click="LoginButtonPressed" :disabled = "Username.length==0 || Username.length < 3 || Username.length>30 || Username.trim().length<3 || Username.trim().length>30" id="LoginButton" class="btn btn-dark" style="margin-top: 1em; margin-bottom: 1em;" >
				<i class="fa-solid fa-right-to-bracket"> Login/Register </i>
			
			</button>	
		</div>
		<span class="fa-solid"></span>

		



		
	</div>
	
</template>

<style>
	#LoginPageContainer{
		background: rgb(87,32,122);
		background: -moz-linear-gradient(68deg, rgba(87,32,122,1) 15%, rgba(104,20,138,1) 50%, rgba(87,32,122,1) 85%);
		background: -webkit-linear-gradient(68deg, rgba(87,32,122,1) 15%, rgba(104,20,138,1) 50%, rgba(87,32,122,1) 85%);
		background: linear-gradient(68deg, rgba(87,32,122,1) 15%, rgba(104,20,138,1) 50%, rgba(87,32,122,1) 85%);
		filter: progid:DXImageTransform.Microsoft.gradient(startColorstr="#57207a",endColorstr="#57207a",GradientType=1); 
	}
	#LoginButton:hover{
		transform: scale(1.1);
	}

</style>
